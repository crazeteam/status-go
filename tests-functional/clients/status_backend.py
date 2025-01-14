import json
import logging
import time
import random
import threading
import requests
from tenacity import retry, stop_after_delay, wait_fixed

from clients.signals import SignalClient
from clients.rpc import RpcClient
from datetime import datetime
from conftest import option
from constants import user_1, DEFAULT_DISPLAY_NAME


class StatusBackend(RpcClient, SignalClient):

    def __init__(self, await_signals=[], url=None):
        try:
            url = url if url else random.choice(option.status_backend_urls)
        except IndexError:
            raise Exception("Not enough status-backend containers, please add more")
        option.status_backend_urls.remove(url)

        self.api_url = f"{url}/statusgo"
        self.ws_url = f"{url}".replace("http", "ws")
        self.rpc_url = f"{url}/statusgo/CallRPC"

        RpcClient.__init__(self, self.rpc_url)
        SignalClient.__init__(self, self.ws_url, await_signals)

        websocket_thread = threading.Thread(target=self._connect)
        websocket_thread.daemon = True
        websocket_thread.start()

    def api_request(self, method, data, url=None):
        url = url if url else self.api_url
        url = f"{url}/{method}"
        logging.info(f"Sending POST request to url {url} with data: {json.dumps(data, sort_keys=True, indent=4)}")
        response = requests.post(url, json=data)
        logging.info(f"Got response: {response.content}")
        return response

    def verify_is_valid_api_response(self, response):
        assert response.status_code == 200, f"Got response {response.content}, status code {response.status_code}"
        assert response.content
        logging.info(f"Got response: {response.content}")
        try:
            assert not response.json()["error"]
        except json.JSONDecodeError:
            raise AssertionError(
                f"Invalid JSON in response: {response.content}")
        except KeyError:
            pass

    def api_valid_request(self, method, data):
        response = self.api_request(method, data)
        self.verify_is_valid_api_response(response)
        return response

    def init_status_backend(self, data_dir="/"):
        method = "InitializeApplication"
        data = {
            "dataDir": data_dir,
            "logEnabled": True,
            "logLevel": "DEBUG",
            "apiLogging": True,
        }
        return self.api_valid_request(method, data)

    def create_account_and_login(self, data_dir="/", display_name=DEFAULT_DISPLAY_NAME, password=user_1.password):
        method = "CreateAccountAndLogin"
        data = {
            "rootDataDir": data_dir,
            "kdfIterations": 256000,
            "displayName": display_name,
            "password": password,
            "customizationColor": "primary",
            "logEnabled": True,
            "logLevel": "DEBUG",
        }
        return self.api_valid_request(method, data)

    def restore_account_and_login(self, data_dir="/",display_name=DEFAULT_DISPLAY_NAME, user=user_1,
                                  network_id=31337):
        method = "RestoreAccountAndLogin"
        data = {
            "rootDataDir": data_dir,
            "kdfIterations": 256000,
            "displayName": display_name,
            "password": user.password,
            "mnemonic": user.passphrase,
            "customizationColor": "blue",
            "logEnabled": True,
            "logLevel": "DEBUG",
            "testNetworksEnabled": False,
            "networkId": network_id,
            "networksOverride": [
                {
                    "ChainID": network_id,
                    "ChainName": "Anvil",
                    "DefaultRPCURL": "http://anvil:8545",
                    "RPCURL": "http://anvil:8545",
                    "ShortName": "eth",
                    "NativeCurrencyName": "Ether",
                    "NativeCurrencySymbol": "ETH",
                    "NativeCurrencyDecimals": 18,
                    "IsTest": False,
                    "Layer": 1,
                    "Enabled": True
                }
            ]
        }
        return self.api_valid_request(method, data)

    def login(self, keyUid, user=user_1):
        method = "LoginAccount"
        data = {
            "password": user.password,
            "keyUid": keyUid,
            "kdfIterations": 256000,
        }
        return self.api_valid_request(method, data)

    def logout(self, user=user_1):
        method = "Logout"
        return self.api_valid_request(method, {})

    def restore_account_and_wait_for_rpc_client_to_start(self, timeout=60):
        self.restore_account_and_login()
        start_time = time.time()
        # ToDo: change this part for waiting for `node.login` signal when websockets are migrated to StatusBackend
        while time.time() - start_time <= timeout:
            try:
                self.rpc_valid_request(method='accounts_getKeypairs')
                return
            except AssertionError:
                time.sleep(3)
        raise TimeoutError(f"RPC client was not started after {timeout} seconds")

    @retry(stop=stop_after_delay(10), wait=wait_fixed(0.5), reraise=True)
    def start_messenger(self, params=[]):
        method = "wakuext_startMessenger"
        response = self.rpc_request(method, params)
        json_response = response.json()

        if 'error' in json_response:
            assert json_response['error']['code'] == -32000
            assert json_response['error']['message'] == "messenger already started"
            return

        self.verify_is_valid_json_rpc_response(response)

    def start_wallet(self, params=[]):
        method = "wallet_startWallet"
        response = self.rpc_request(method, params)
        self.verify_is_valid_json_rpc_response(response)

    def get_settings(self, params=[]):
        method = "settings_getSettings"
        response = self.rpc_request(method, params)
        self.verify_is_valid_json_rpc_response(response)

    def get_accounts(self, params=[]):
        method = "accounts_getAccounts"
        response = self.rpc_request(method, params)
        self.verify_is_valid_json_rpc_response(response)
        return response.json()

    def get_pubkey(self, display_name):
        response = self.get_accounts()
        accounts = response.get("result", [])
        for account in accounts:
            if account.get("name") == display_name:
                return account.get("public-key")
        raise ValueError(f"Public key not found for display name: {display_name}")

    def send_contact_request(self, params=[]):
        method = "wakuext_sendContactRequest"
        response = self.rpc_request(method, params)
        self.verify_is_valid_json_rpc_response(response)
        return response.json()

    def send_message(self, params=[]):
        method = "wakuext_sendOneToOneMessage"
        response = self.rpc_request(method, params)
        self.verify_is_valid_json_rpc_response(response)
        return response.json()
