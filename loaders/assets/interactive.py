#!/usr/bin/python3

from ibapi.client import EClient
from ibapi.wrapper import EWrapper
from ibapi.contract import Contract

from ibapi.account_summary_tags import *

import requests
import json
from datetime import datetime

import threading
import time

class IBapi(EWrapper, EClient):
    def __init__(self):
        EClient.__init__(self, self)

    def accountSummary(self, reqId: int, account: str, tag: str, value: str, currency: str):
        super().accountSummary(reqId, account, tag, value, currency)
        print("AccountSummary. ReqId:", reqId, "Account:", account,
                "Tag: ", tag, "Value:", value, "Currency:", currency)
        if tag == "TotalCashValue":
            data = {"amount": value,
                    "assetId":13,
                    "date": datetime.today().strftime('%Y-%m-%d')}
            print(data)
            r = requests.post(url = "https://debian.home:8000/assets/series", data=json.dumps(data), verify=False)
            print(r)

    def position(self, account: str, contract: Contract, position: float,
            avgCost: float):
        super().position(account, contract, position, avgCost)
        print("Position.", "Account:", account, "Symbol:", contract.symbol, "SecType:",
                contract.secType, "Currency:", contract.currency,
                "Position:", position, "Avg cost:", avgCost)


def run_loop():
    app.run()

app = IBapi()
app.connect("127.0.0.1", 7496, clientId=11)

api_thread = threading.Thread(target=run_loop, daemon=True)
api_thread.start()

time.sleep(2)
app.reqAccountSummary(9001, 'All', AccountSummaryTags.AllTags)
app.reqPositions()

