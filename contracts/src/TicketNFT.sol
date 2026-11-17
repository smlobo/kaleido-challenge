// SPDX-License-Identifier: MIT
pragma solidity 0.8.38;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";

/// @title TicketNFT
/// @notice One contract per event. Each token is one ticket. The box office
/// (MINTER_ROLE) mints, holders can list a ticket for resale at no more than
/// the event's resale cap, anyone can buy a listed ticket, and the gate
/// (GATE_ROLE) burns a ticket at check-in.
///
/// Tickets cannot be moved with the standard ERC-721 transfer and approval
/// functions; buyResale is the only way a ticket changes hands, so the resale
/// cap cannot be bypassed.
///
/// No money moves yet: prices are checked and recorded, and payment is
/// collected off-chain. On-chain payment (TicketUSD) is a later milestone.
contract TicketNFT is ERC721, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant GATE_ROLE = keccak256("GATE_ROLE");

    /// @notice Most tickets that can ever be minted; burned tickets still count.
    uint32 public immutable maxSupply;
    /// @notice Resale cap in basis points of face value (11000 = 110%).
    uint32 public immutable maxResaleBps;
    /// @notice Event start, unix seconds. Minting, listing and resale stop here.
    uint64 public immutable eventStart;

    /// @notice Tickets minted so far, including ones since checked in.
    uint32 public mintedCount;

    struct Ticket {
        uint256 faceValue;
        bytes32 recordHash; // keccak256 of the off-chain ticket record
        uint256 listingPrice; // 0 = not listed
    }

    mapping(uint256 tokenId => Ticket) private _tickets;
    string private _base;

    event TicketMinted(uint256 indexed tokenId, address indexed to, uint256 faceValue, bytes32 recordHash);
    event TicketListed(uint256 indexed tokenId, address indexed seller, uint256 price);
    event ListingCancelled(uint256 indexed tokenId);
    event TicketResold(uint256 indexed tokenId, address indexed from, address indexed to, uint256 price);
    event CheckedIn(uint256 indexed tokenId, address indexed holder);

    error InvalidConfig();
    error EventStarted();
    error SoldOut();
    error EmptyRecordHash();
    error NotTicketHolder(uint256 tokenId);
    error ZeroPrice();
    error PriceAboveCap(uint256 price, uint256 cap);
    error NotListed(uint256 tokenId);
    error PriceAboveMaxPrice(uint256 price, uint256 maxPrice);
    error BuyerIsSeller();
    error TransfersDisabled();

    /// @param admin   gets DEFAULT_ADMIN_ROLE (the platform deployer)
    /// @param minter  gets MINTER_ROLE (the organizer's box office key)
    /// @param gate    gets GATE_ROLE (the organizer's gate key)
    constructor(
        string memory name_,
        string memory symbol_,
        uint32 maxSupply_,
        uint32 maxResaleBps_,
        uint64 eventStart_,
        string memory baseURI_,
        address admin,
        address minter,
        address gate
    ) ERC721(name_, symbol_) {
        // Block timestamps can drift by seconds; that is fine for an event start.
        // forge-lint: disable-next-line(block-timestamp)
        if (maxSupply_ == 0 || eventStart_ <= block.timestamp) revert InvalidConfig();
        if (admin == address(0) || minter == address(0) || gate == address(0)) revert InvalidConfig();
        maxSupply = maxSupply_;
        maxResaleBps = maxResaleBps_;
        eventStart = eventStart_;
        _base = baseURI_;
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(MINTER_ROLE, minter);
        _grantRole(GATE_ROLE, gate);
    }

    modifier beforeStart() {
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= eventStart) revert EventStarted();
        _;
    }

    /// @notice Mint a ticket. The service picks tokenId (the ticket's serial
    /// number within the event), so a retried mint reverts instead of creating
    /// a duplicate.
    function mintTicket(address to, uint256 tokenId, uint256 faceValue, bytes32 recordHash)
        external
        onlyRole(MINTER_ROLE)
        beforeStart
    {
        if (mintedCount >= maxSupply) revert SoldOut();
        if (recordHash == bytes32(0)) revert EmptyRecordHash();
        mintedCount++;
        _tickets[tokenId] = Ticket({faceValue: faceValue, recordHash: recordHash, listingPrice: 0});
        // _mint, not _safeMint: tickets go to the service's custodial wallets
        // (plain accounts), and _safeMint would call into the recipient if it
        // were a contract. Reverts if tokenId already exists.
        // forge-lint: disable-next-line(unsafe-oz-erc721-mint)
        _mint(to, tokenId);
        emit TicketMinted(tokenId, to, faceValue, recordHash);
    }

    /// @notice List (or re-price) your ticket for resale.
    function listForResale(uint256 tokenId, uint256 price) external beforeStart {
        _requireHolder(tokenId);
        if (price == 0) revert ZeroPrice();
        uint256 cap = resaleCap(tokenId);
        if (price > cap) revert PriceAboveCap(price, cap);
        _tickets[tokenId].listingPrice = price;
        // No external call happens here; the linter flags the view helpers.
        // forge-lint: disable-next-line(reentrancy-events)
        emit TicketListed(tokenId, msg.sender, price);
    }

    /// @notice Take your ticket off the resale market.
    function cancelListing(uint256 tokenId) external {
        _requireHolder(tokenId);
        if (_tickets[tokenId].listingPrice == 0) revert NotListed(tokenId);
        _tickets[tokenId].listingPrice = 0;
        emit ListingCancelled(tokenId);
    }

    /// @notice Buy a listed ticket. maxPrice protects the buyer: if the seller
    /// re-lists at a higher price after the buyer looked, this reverts.
    function buyResale(uint256 tokenId, uint256 maxPrice) external beforeStart {
        uint256 price = _tickets[tokenId].listingPrice;
        if (price == 0) revert NotListed(tokenId);
        if (price > maxPrice) revert PriceAboveMaxPrice(price, maxPrice);
        address seller = ownerOf(tokenId);
        if (msg.sender == seller) revert BuyerIsSeller();
        _tickets[tokenId].listingPrice = 0;
        _transfer(seller, msg.sender, tokenId);
        emit TicketResold(tokenId, seller, msg.sender, price);
    }

    /// @notice Check a ticket in at the door: burns it and clears any listing.
    function checkIn(uint256 tokenId) external onlyRole(GATE_ROLE) {
        address holder = ownerOf(tokenId); // reverts for unknown or already checked-in tickets
        _tickets[tokenId].listingPrice = 0;
        _burn(tokenId);
        emit CheckedIn(tokenId, holder);
    }

    /// @notice Face value, record hash and current listing price (0 = not
    /// listed). Face value and record hash stay readable after check-in.
    function ticketInfo(uint256 tokenId)
        external
        view
        returns (uint256 faceValue, bytes32 recordHash, uint256 listingPrice)
    {
        Ticket storage t = _tickets[tokenId];
        return (t.faceValue, t.recordHash, t.listingPrice);
    }

    /// @notice Highest price this ticket may be listed at.
    function resaleCap(uint256 tokenId) public view returns (uint256) {
        return Math.mulDiv(_tickets[tokenId].faceValue, maxResaleBps, 10_000);
    }

    // ---- Standard transfers and approvals are disabled ----

    function transferFrom(address, address, uint256) public pure override {
        revert TransfersDisabled();
    }

    /// @dev The 3-argument safeTransferFrom calls this one, so both are blocked.
    function safeTransferFrom(address, address, uint256, bytes memory) public pure override {
        revert TransfersDisabled();
    }

    function approve(address, uint256) public pure override {
        revert TransfersDisabled();
    }

    function setApprovalForAll(address, bool) public pure override {
        revert TransfersDisabled();
    }

    // ---- Internals ----

    function _requireHolder(uint256 tokenId) private view {
        if (ownerOf(tokenId) != msg.sender) revert NotTicketHolder(tokenId);
    }

    function _baseURI() internal view override returns (string memory) {
        return _base;
    }

    function supportsInterface(bytes4 interfaceId) public view override(ERC721, AccessControl) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
