'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');

const ccpPath = path.resolve(__dirname, 'connection-org1.json');

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('user1');
        if (!identity) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('channel1');

        // Get the contract from the network.
        const contract = network.getContract('mycc');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        await contract.submitTransaction("invoke", "request_access", "200", "mickey@disney.com", "eyJyZXFfdHlwZSI6ImVtcGxveWVyIiwicmVxX25hbWUiOiJNaWNrZXkgTW91c2UiLCJyZXFfY29tcGFueSI6IkRpc25leSIsInJlcV9qb2JfcG9zaXRpb24iOiIiLCJhY2NlcHRlZCI6ZmFsc2UsImNpdHkiOiJUb2tpbyIsInN0YXRlIjoiVG9raW8iLCJjcmVhdGVkX2RhdGUiOiIyMDA5LTExLTEwVDIzOjAwOjAwWiIsImFjY2VwdGVkX2RhdGUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImRuYV90eXBlcyI6WzEsMiwzXX0=");
        await contract.submitTransaction("invoke", "request_accept", "9", "200", "WyJkbmFfdHlwZV8xIiwgImRuYV90eXBlXzIiXQ==", "mickey@disney.com");
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
