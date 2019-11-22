package example

import util.movieYearRegex;
import org.apache.spark._
import SparkContext._
import org.apache.spark.sql.SparkSession
import scala.util.control.Breaks
import scala.collection.JavaConversions._
import scala.collection.JavaConverters._
import com.hankcs.hanlp.HanLP
import org.apache.spark.sql.Row
import scala.collection.mutable
import scala.collection.mutable.ListBuffer
import org.apache.spark.sql.types.{DoubleType, StringType, StructField, StructType}

object Movies {
  def main(args: Array[String]) {
    val sparkSession = SparkSession.builder().appName("base-content-Recommand").getOrCreate()

    val rdd = sparkSession.sparkContext.textFile("file:///opt/spark/ratings.dat").map(_.split("\\::"))
    //(movie,avg_rate)
    val movieAvgRate = rdd.map {
      f =>
        (f(1), f(2))
    }
    //电影id，名称，以及genre类别
    val moviesData = sparkSession.sparkContext.textFile("file:///opt/spark/movies.dat").map(_.split("\\::"))
    //电影tag
    val tagsData = sparkSession.sparkContext.textFile("file:///opt/spark/tags.dat").map(_.split("\\::"))
    val tagsStandardize = tagsData.map {
      f =>
        (f(1), f(2))
    }
    val tagsStandardizeTmp = tagsStandardize.collect()
    val tagsSimi = tagsStandardize.map {
      f =>
        var retTag = f._2
        if (f._2.toString.split(" ").size == 1) {
          var simiTmp = ""
          val tagsTmpStand = tagsStandardizeTmp
            .filter(_._2.toString.split(" ").size != 1)
            .filter(f._2.toString.size < _._2.toString.size)
            .sortBy(_._2.toString.size)
          var x = 0
          val loop = new Breaks
          tagsTmpStand.map {
            tagTmp =>
              val flag = getEditSize(f._2.toString, tagTmp._2.toString)
              if (flag == 1) {
                retTag = tagTmp._2
                loop.break()
              }
          }
          ((f._1, retTag), 1)
        } else {
          ((f._1, f._2), 1)
        }
    }

    val movieTag = tagsSimi.reduceByKey(_ + _).groupBy(k => k._1._1).map {
      f =>
        (f._1, f._2.map {
          ff =>
            (ff._1._2, ff._2)
        }.toList.sortBy(_._2).reverse.take(10).toMap)
    }

    val moviesGenresTitleYear = moviesData.map {
      f =>
        val movieid = f(0)
        val title = f(1)
        val genres = f(2).toString.split("|").toList.take(10)
        val titleWorlds = title.toString().split(" ").toList
        val year = movieYearRegex.movieYearReg(title.toString)
        (movieid, (genres, titleWorlds, year))
    }
    val movieContent = movieTag.join(movieAvgRate).join(moviesGenresTitleYear).map {
      f =>
        //(movie,tagList,titleList,year,genreList,rate)
        (f._1, f._2._1._1, f._2._2._2, f._2._2._3, f._2._2._1, f._2._1._2)
    }

    val movieContentTmp = movieContent.filter(f => BigDecimal(f._6).doubleValue() < 3.5).collect()
    val movieContentBase = movieContent.map {
      f =>
        val currentMoiveId = f._1
        val currentTagList = f._2 //[(tag,score)]
        val currentTitleWorldList = f._3
        val currentYear = f._4
        val currentGenreList = f._5
        val currentRate = BigDecimal(f._6).doubleValue()

        val recommandMovies = movieContentTmp.map {
          ff =>
            val tagSimi = getCosTags(currentTagList, ff._2)
            val titleSimi = getCosList(currentTitleWorldList, ff._3)
            val genreSimi = getCosList(currentGenreList, ff._5)
            val yearSimi = getYearSimi(currentYear, ff._4)
            val rateSimi = getRateSimi(BigDecimal(ff._6).doubleValue())
            val score = 0.4 * genreSimi + 0.25 * tagSimi + 0.1 * yearSimi + 0.05 * titleSimi + 0.2 * rateSimi
            (ff._1, score)
        }.toList.sortBy(k => k._2).reverse.take(20)
        (currentMoiveId, recommandMovies)
    }.flatMap(f => f._2.map(k => (f._1, k._1, k._2))).map(f => Row(f._1, f._2, f._3))
    val schemaString2 = "movieid movieid_re score"
    val schemaContentBase = StructType(schemaString2.split(" ")
      .map(fieldName=>StructField(fieldName,if (fieldName.equals("score")) DoubleType else  StringType,true)))
    val movieContentBaseDataFrame = sparkSession.createDataFrame(movieContentBase, schemaContentBase)
	//dataFrame
	
    movieContentBaseDataFrame.write.csv("file:///opt/spark/result.csv")
  }

  def getRateSimi(rate2: Double): Double = {
    if (rate2 >= 5) {
      1
    } else {
      rate2 / 5
    }
  }

  def getYearSimi(year1: Int, year2: Int): Double = {
    val count = Math.abs(year1 - year2)
    if (count > 10) {
      0
    } else {
      (1 - count) / 10
    }
  }
  
  def getCosList(listTags1: List[String], listTags2: List[String]): Double = {

    var xySum: Double = 0
    var aSquareSum: Double = 0
    var bSquareSum: Double = 0

    listTags1.union(listTags2).map {
      f =>
        if (listTags1.contains(f)) aSquareSum += 1
        if (listTags2.contains(f)) bSquareSum += 1
        if (listTags1.contains(f) && listTags2.contains(f)) xySum += 1
    }

    if (aSquareSum != 0 && bSquareSum != 0) {
      xySum / (Math.sqrt(aSquareSum) * Math.sqrt(bSquareSum))
    } else {
      0
    }

  }

  def getCosTags(listTagsCurrent: Map[String, Int], listTagsTmp: Map[String, Int]): Double = {

    var xySum: Double = 0
    var aSquareSum: Double = 0
    var bSquareSum: Double = 0

    val tagsA = listTagsCurrent.map(f => f._1).toList
    val tagsB = listTagsTmp.map(f => f._1).toList
    tagsA.union(tagsB).map {
      f =>
        if (listTagsCurrent.contains(f)) (aSquareSum += listTagsCurrent.get(f).get * listTagsCurrent.get(f).get)
        if (listTagsTmp.contains(f)) (bSquareSum += listTagsTmp.get(f).get * listTagsTmp.get(f).get)
        if (listTagsCurrent.contains(f) && listTagsTmp.contains(f)) (xySum += listTagsCurrent.get(f).get * listTagsTmp.get(f).get)
    }

    if (aSquareSum != 0 && bSquareSum != 0) {
      xySum / (Math.sqrt(aSquareSum) * Math.sqrt(bSquareSum))
    } else {
      0
    }

  }

  def getEditSize(str1: String, str2: String): Int = {
    if (str2.size > str1.size) {
      0
    } else {
      var count = 0
      val loop = new Breaks
      val lengthStr2 = str2.getBytes().length
      var i = 0
      for (i <- 1 to lengthStr2) {
        if (str2.getBytes()(i) == str1.getBytes()(i)) {
          count += 1
        } else {
          loop.break()
        }
      }
      if (count.asInstanceOf[Double] / str1.getBytes().size.asInstanceOf[Double] >= (1 - 0.286)) {
        1
      } else {
        0
      }
    }
  }

  /**
   * 向量的模长
   * @param vec
   */
  def module(vec: Vector[Double]): Double = {
    // math.sqrt( vec.map(x=>x*x).sum )
    math.sqrt(vec.map(math.pow(_, 2)).sum)
  }

}