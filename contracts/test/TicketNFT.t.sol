// SPDX-License-Identifier: MIT
pragma solidity 0.8.38;

import {Test} from "forge-std/Test.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IERC721Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {TicketNFT} from "../src/TicketNFT.sol";

contract TicketNFTTest is Test {
    TicketNFT ticketnft;

    address admin = makeAddr("admin");
    address boxOffice = makeAddr("boxOffice");
    address gate = makeAddr("gate");
    address alice = makeAddr("alice");
    address bob = makeAddr("bob");
    address carol = makeAddr("carol");

    uint32 constant SUPPLY = 3;
    uint32 constant RESALE_BPS = 11_000; // 110%
    uint256 constant FACE = 100e6; // 100 units with 6 decimals
    uint256 constant CAP = 110e6;
    bytes32 constant RECORD = keccak256("ticket #1 record");
    string constant BASE_URL = "https://tickets.example/meta/";
    uint64 start;

    event TicketMinted(uint256 indexed tokenId, address indexed to, uint256 faceValue, bytes32 recordHash);
    event TicketListed(uint256 indexed tokenId, address indexed seller, uint256 price);
    event ListingCancelled(uint256 indexed tokenId);
    event TicketResold(uint256 indexed tokenId, address indexed from, address indexed to, uint256 price);
    event CheckedIn(uint256 indexed tokenId, address indexed holder);

    function setUp() public {
        // Set to Jan 15, 2027
        vm.warp(1_800_000_000);
        start = uint64(block.timestamp + 7 days);
        ticketnft = _deploy(SUPPLY, start);
    }

    function _deploy(uint32 supply, uint64 eventStart) internal returns (TicketNFT) {
        return new TicketNFT("Rock Fest", "ROCK", supply, RESALE_BPS, eventStart, BASE_URL, admin, boxOffice, gate);
    }

    function _mintAsBoxOffice(address to, uint256 tokenId) internal {
        vm.prank(boxOffice);
        ticketnft.mintTicket(to, tokenId, FACE, RECORD);
    }

    function _listAs(address seller, uint256 tokenId, uint256 price) internal {
        vm.prank(seller);
        ticketnft.listForResale(tokenId, price);
    }

    // ---- Deploy ----

    function test_Constructor_SetsRulesAndRoles() public view {
        assertEq(ticketnft.name(), "Rock Fest");
        assertEq(ticketnft.symbol(), "ROCK");
        assertEq(ticketnft.maxSupply(), SUPPLY);
        assertEq(ticketnft.maxResaleBps(), RESALE_BPS);
        assertEq(ticketnft.eventStart(), start);
        assertTrue(ticketnft.hasRole(ticketnft.DEFAULT_ADMIN_ROLE(), admin));
        assertTrue(ticketnft.hasRole(ticketnft.MINTER_ROLE(), boxOffice));
        assertTrue(ticketnft.hasRole(ticketnft.GATE_ROLE(), gate));
        assertFalse(ticketnft.hasRole(ticketnft.MINTER_ROLE(), admin), "admin is not a minter by default");
    }

    function test_Constructor_RejectsBadConfig() public {
        vm.expectRevert(TicketNFT.InvalidConfig.selector);
        _deploy(0, start);

        vm.expectRevert(TicketNFT.InvalidConfig.selector);
        _deploy(SUPPLY, uint64(block.timestamp)); // start must be in the future

        vm.expectRevert(TicketNFT.InvalidConfig.selector);
        new TicketNFT("x", "X", 1, RESALE_BPS, start, "", admin, address(0), gate);
    }

    function test_SupportsInterfaces() public view {
        // The compiler XORs IERC721's own function selectors (not inherited
        // IERC165's), which is the ID the ERC-721 standard publishes.
        assertEq(type(IERC721).interfaceId, bytes4(0x80ac58cd));
        assertTrue(ticketnft.supportsInterface(type(IERC721).interfaceId), "ERC-721");
        assertTrue(ticketnft.supportsInterface(type(IAccessControl).interfaceId), "AccessControl");
    }

    // ---- Mint ----

    function test_Mint_StoresTicketAndEmits() public {
        vm.expectEmit(address(ticketnft));
        emit TicketMinted(1, alice, FACE, RECORD);
        _mintAsBoxOffice(alice, 1);

        assertEq(ticketnft.ownerOf(1), alice);
        assertEq(ticketnft.mintedCount(), 1);
        assertEq(ticketnft.tokenURI(1), string.concat(BASE_URL, "1"));
        (uint256 face, bytes32 record, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(face, FACE);
        assertEq(record, RECORD);
        assertEq(listing, 0);
    }

    function test_Mint_OnlyMinter() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, alice, ticketnft.MINTER_ROLE())
        );
        vm.prank(alice);
        ticketnft.mintTicket(alice, 1, FACE, RECORD);
    }

    function test_Mint_DuplicateTokenIdReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721InvalidSender.selector, address(0)));
        _mintAsBoxOffice(bob, 1); // a retried mint of the same ticket
    }

    function test_Mint_StopsAtMaxSupply() public {
        for (uint256 id = 1; id <= SUPPLY; id++) {
            _mintAsBoxOffice(alice, id);
        }
        vm.expectRevert(TicketNFT.SoldOut.selector);
        _mintAsBoxOffice(alice, 99);
    }

    function test_Mint_CheckedInTicketsStillCountTowardSupply() public {
        for (uint256 id = 1; id <= SUPPLY; id++) {
            _mintAsBoxOffice(alice, id);
        }
        vm.prank(gate);
        ticketnft.checkIn(1);
        vm.expectRevert(TicketNFT.SoldOut.selector);
        _mintAsBoxOffice(alice, 99);
    }

    function test_Mint_RejectsEmptyRecordHash() public {
        vm.expectRevert(TicketNFT.EmptyRecordHash.selector);
        vm.prank(boxOffice);
        ticketnft.mintTicket(alice, 1, FACE, bytes32(0));
    }

    function test_Mint_StopsAtEventStart() public {
        vm.warp(start); // exactly the start time counts as started
        vm.expectRevert(TicketNFT.EventStarted.selector);
        _mintAsBoxOffice(alice, 1);
    }

    // ---- List ----

    function test_List_AtExactlyTheCap() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectEmit(address(ticketnft));
        emit TicketListed(1, alice, CAP);
        _listAs(alice, 1, CAP);
        (,, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(listing, CAP);
        assertEq(ticketnft.resaleCap(1), CAP);
    }

    function test_List_OneAboveTheCapReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.PriceAboveCap.selector, CAP + 1, CAP));
        _listAs(alice, 1, CAP + 1);
    }

    function testFuzz_List_CapBoundary(uint256 price) public {
        _mintAsBoxOffice(alice, 1);
        price = bound(price, 1, CAP * 3);
        if (price > CAP) {
            vm.expectRevert(abi.encodeWithSelector(TicketNFT.PriceAboveCap.selector, price, CAP));
        }
        _listAs(alice, 1, price);
    }

    function test_List_ZeroPriceReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(TicketNFT.ZeroPrice.selector);
        _listAs(alice, 1, 0);
    }

    function test_List_OnlyHolder() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.NotTicketHolder.selector, 1));
        _listAs(bob, 1, FACE);
    }

    function test_List_StopsAtEventStart() public {
        _mintAsBoxOffice(alice, 1);
        vm.warp(start);
        vm.expectRevert(TicketNFT.EventStarted.selector);
        _listAs(alice, 1, FACE);
    }

    function test_List_CanReprice() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        _listAs(alice, 1, CAP);
        (,, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(listing, CAP);
    }

    // ---- Cancel ----

    function test_Cancel_ClearsListing() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        vm.expectEmit(address(ticketnft));
        emit ListingCancelled(1);
        vm.prank(alice);
        ticketnft.cancelListing(1);
        (,, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(listing, 0);
    }

    function test_Cancel_NotListedReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.NotListed.selector, 1));
        vm.prank(alice);
        ticketnft.cancelListing(1);
    }

    function test_Cancel_OnlyHolder() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.NotTicketHolder.selector, 1));
        vm.prank(bob);
        ticketnft.cancelListing(1);
    }

    // ---- Buy ----

    function test_Buy_MovesTicketAndClearsListing() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, CAP);
        vm.expectEmit(address(ticketnft));
        emit TicketResold(1, alice, bob, CAP);
        vm.prank(bob);
        ticketnft.buyResale(1, CAP);

        assertEq(ticketnft.ownerOf(1), bob);
        assertEq(ticketnft.balanceOf(alice), 0);
        (uint256 face,, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(listing, 0);
        assertEq(face, FACE, "face value survives resale, so the cap stays the same");
    }

    function test_Buy_NewHolderCanResellUnderTheSameCap() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, CAP);
        vm.prank(bob);
        ticketnft.buyResale(1, CAP);

        _listAs(bob, 1, CAP);
        vm.prank(carol);
        ticketnft.buyResale(1, CAP);
        assertEq(ticketnft.ownerOf(1), carol);
    }

    function test_Buy_NotListedReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.NotListed.selector, 1));
        vm.prank(bob);
        ticketnft.buyResale(1, CAP);
    }

    function test_Buy_SellerRaisedPriceReverts() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE); // bob sees 100
        _listAs(alice, 1, CAP); // alice re-lists at 110 before bob buys
        vm.expectRevert(abi.encodeWithSelector(TicketNFT.PriceAboveMaxPrice.selector, CAP, FACE));
        vm.prank(bob);
        ticketnft.buyResale(1, FACE);
    }

    function test_Buy_FromYourselfReverts() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        vm.expectRevert(TicketNFT.BuyerIsSeller.selector);
        vm.prank(alice);
        ticketnft.buyResale(1, FACE);
    }

    function test_Buy_StopsAtEventStart() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        vm.warp(start);
        vm.expectRevert(TicketNFT.EventStarted.selector);
        vm.prank(bob);
        ticketnft.buyResale(1, FACE);
    }

    // ---- Check-in ----

    function test_CheckIn_BurnsAndKeepsRecord() public {
        _mintAsBoxOffice(alice, 1);
        _listAs(alice, 1, FACE);
        vm.warp(start + 1 hours); // at the door, after the start
        vm.expectEmit(address(ticketnft));
        emit CheckedIn(1, alice);
        vm.prank(gate);
        ticketnft.checkIn(1);

        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, 1));
        ticketnft.ownerOf(1);
        (uint256 face, bytes32 record, uint256 listing) = ticketnft.ticketInfo(1);
        assertEq(face, FACE);
        assertEq(record, RECORD);
        assertEq(listing, 0, "listing cleared on check-in");
        assertEq(ticketnft.balanceOf(alice), 0);
    }

    function test_CheckIn_OnlyGate() public {
        _mintAsBoxOffice(alice, 1);
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, boxOffice, ticketnft.GATE_ROLE())
        );
        vm.prank(boxOffice);
        ticketnft.checkIn(1);
    }

    function test_CheckIn_TwiceReverts() public {
        _mintAsBoxOffice(alice, 1);
        vm.startPrank(gate);
        ticketnft.checkIn(1);
        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, 1));
        ticketnft.checkIn(1);
        vm.stopPrank();
    }

    // ---- Standard transfers are blocked ----

    function test_Transfers_AllBlocked() public {
        _mintAsBoxOffice(alice, 1);
        vm.startPrank(alice);
        vm.expectRevert(TicketNFT.TransfersDisabled.selector);
        ticketnft.transferFrom(alice, bob, 1);
        vm.expectRevert(TicketNFT.TransfersDisabled.selector);
        ticketnft.safeTransferFrom(alice, bob, 1);
        vm.expectRevert(TicketNFT.TransfersDisabled.selector);
        ticketnft.safeTransferFrom(alice, bob, 1, "");
        vm.expectRevert(TicketNFT.TransfersDisabled.selector);
        ticketnft.approve(bob, 1);
        vm.expectRevert(TicketNFT.TransfersDisabled.selector);
        ticketnft.setApprovalForAll(bob, true);
        vm.stopPrank();
        assertEq(ticketnft.ownerOf(1), alice);
    }
}
