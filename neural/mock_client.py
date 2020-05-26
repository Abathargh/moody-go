from xmlrpc.client import ServerProxy

if __name__ == "__main__":
    datatypes = ["situation", "audio"]
    dataset = "test.csv"
    data = [0]

    with ServerProxy("http://localhost:9998/neural") as proxy:
        try:
            res = proxy.train(dataset, datatypes)
            print(res)
        except Exception as e:
            print(e)
