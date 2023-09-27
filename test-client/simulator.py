import sys
import requests
import pandas as pd

def simulate(epoch):
    keyList = ["distPoolId", "prizeDenomId", "name", "type", "amount", "odds"]
    result = {}

    for _ in range(epoch):
        body = {
            'addr': 'cre1zakl8nyee0ln4lleupt7c9eex7tg2urw784r9j',
            'promotionId': 49,
            'gameId': 9,
            'usedTicketQty': 1
        }
        res = requests.post('http://127.0.0.1:8080/game-mgmt/start', json=body)

        resData = res.json()['data']

        temp = {}
        for key in keyList:
            temp[key] = resData[key]

        key = resData['prizeId']
        if key in result:
            result[key]['count']+=1
        else:
            result[key] = temp
            result[key]['count'] = 1

    printSimulateResult(result, epoch)

def printSimulateResult(result, epoch):
    headers = ["Prize ID", "Prize Info", "등록된 확률", "당첨 횟수", "계산된 확률"]
    excel = {}
    for header in headers:
        excel[header] = []

    for k, v in sorted(result.items()):
        excel[headers[0]].append(k)
        excel[headers[1]].append("{} {} {}".format(v['name'], v['type'], v['amount']))
        excel[headers[2]].append("{}%".format(v['odds']*100))
        excel[headers[3]].append( v['count'] )
        excel[headers[4]].append("{}%".format((v['count']/epoch) * 100))

    raw_data = pd.DataFrame(excel)
    raw_data.to_excel(excel_writer='{}result.xlsx'.format(epoch), index=False)

def main():
    epoch = int(sys.argv[1])
    simulate(epoch)
main()
