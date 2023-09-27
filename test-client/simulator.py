import sys
import requests

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
    for k, v in result.items():
        print('='*20)
        print('prize id:', k)
        print('prize info:', v['name'], v['type'], v['amount'])
        print('registed odds:', v['odds']*100, '%')
        print('win count:', v['count'])
        print('calculated odds:', (v['count']/epoch) * 100, '%')
        print('='*20)

def main():
    epoch = int(sys.argv[1])
    print("# Simulation:", epoch)
    simulate(epoch)
main()
