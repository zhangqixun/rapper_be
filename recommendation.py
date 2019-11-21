from pyspark.ml.evaluation import RegressionEvaluator
from pyspark.ml.recommendation import ALS,ALSModel
from pyspark.sql import Row
from pyspark.sql.session import SparkSession
from pyspark import SparkContext,SparkConf
import pandas as pd
import requests
import json
import pymysql
import sys
from sqlalchemy import create_engine

def update_model():
    url = "http://47.92.86.194:8080/user/footprints/all"
    r = requests.post(url)
    js = json.loads(str(r.content, encoding="utf-8"))
    data = [[int(item["user_id"]), int(item["movie_id"]), int(item["time_on_site"])] for item in js["data"] if int(item["user_id"]) > 0]
    ratings = spark.createDataFrame(data, schema=["user", "item", "rating"])
    # (training,test) = ratings.randomSplit([0.8,0.2])
    als = ALS(maxIter=10, regParam=0.01, userCol="user", itemCol="item",
              ratingCol="rating", coldStartStrategy="drop")
    model = als.fit(ratings)
    # predictions = model.transform(test)
    # evaluator = RegressionEvaluator(metricName="rmse",labelCol="rating",predictionCol="prediction")
    # rmse = evaluator.evaluate(predictions)
    print("-" * 100)
    userRecs = model.recommendForAllUsers(10)  # pyspark.sql.dataframe.DataFrame
    pandas_df = userRecs.toPandas()
    pandas_df["recommendations"] = pandas_df["recommendations"].apply(lambda x: [item[0] for item in x])
    pandas_df["recommendations"] = pandas_df["recommendations"].apply(lambda x: " ".join(list(map(str, x))))
    pandas_df = pandas_df.rename(columns={"recommendations": "recs"})
    pandas_df.to_csv("res.csv", index=False)
    print("finish pandas_df")
    # pandas_df.to_json("res.json",orient="records")

    con = create_engine("mysql+mysqldb://root:123456@47.92.86.194:3306/rapper?charset=utf8")
    pandas_df.to_sql("recs", con=con, if_exists="replace", index=False)

if __name__ == "__main__":
    # if len(sys.argv) != 2:
    #     print("ueage:recommendation.py user_id",file=sys.stderr)
    #     exit(-1)
    appname = "recommendation"
    master = "spark://master:7077"
    conf = SparkConf().setAppName(appname).setMaster(master)
    conf.set("spark.dynamicAllocation.enabled","false") #可选
    sc = SparkContext(conf=conf)
    sc.setLogLevel("ERROR")
    spark = SparkSession(sc)
    # user_id = sys.argv[1]

    # conn = pymysql.connect(
    #     host='47.92.86.194',
    #     port=3306,
    #     user="root", password="123456",
    #     database="rapper",
    #     charset="utf8"
    # )
    # cursor = conn.cursor()
    # effect_row = cursor.execute("select * from recs where user="+user_id)
    # if not cursor.fetchall():
    #     update_model()
    update_model()
    print("finish")

    #use saved model
    # model = ALSModel.load("als_model")
    # # test = spark.createDataFrame([[1,2]],["user","item"])
    # # print("-"*100)
    # # print("-"*100)
    # # predictions = model.transform(test)
    # print("-"*100)
    # pandas_df = model.recommendForAllUsers(100).toPandas()
    # print("-"*100)
    # pandas_df.to_csv("res.csv",index=False)
    # pandas_df = pd.read_csv("res.csv")
    # pat = r"item=(\d*)"
    # mod = re.compile(pat)
    # pandas_df["recommendations"] = pandas_df["recommendations"].apply(lambda x: re.findall(mod, x))
    # pandas_df.to_csv("res.csv", index=False)
    # print("finish")






