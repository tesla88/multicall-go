
    // SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

interface IERC20 {
    function balanceOf(address account) external view returns (uint);
}

contract multiBalance{
    function getBalance(address tokenContract, address[] memory addressArr) public view returns ( uint256[] memory ) {
        uint[] memory returnData = new uint[] (addressArr.length);
        for(uint256 i = 0; i < addressArr.length; i++) {
            returnData[i] = IERC20(tokenContract).balanceOf(addressArr[i]);
        }
        return returnData;
    }
}