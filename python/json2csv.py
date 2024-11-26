import json
import csv

JSON_FILE = "data.json"
CSV_FILE_FORMAT = "data_" # creates "./data_Default.csv" etc.

def json2csv():
    with open(JSON_FILE, "r") as f:
        file = f.read()
        j = json.loads(file)
    for key in j.keys():
        with open(CSV_FILE_FORMAT + key + ".csv" , "w") as c:
            field_names = j[key][0].keys()
            writer = csv.DictWriter(c, field_names)

            writer.writeheader()
            for row in j[key]:
                writer.writerow(row)


json2csv()
