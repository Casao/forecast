package main

import (
  "errors"
  "fmt"
  "forecast/icon"
  "log"
  "os"
  "strings"
  "time"

  "github.com/deckarep/gosx-notifier"
  "github.com/everdev/mack"
  "github.com/getlantern/systray"
  "github.com/mlbright/forecast/v2"
)

var devkey string = os.Getenv("FORECAST_DEVKEY")
var lat string = "33.967754"
var long string = "-84.220302"

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  systray.Run(onReady)
}

func onReady() {
  fore := getForecast(devkey, lat, long)
  ticker := time.NewTicker(time.Minute * 7)
  go func() {
    for _ = range ticker.C {
      checkForecast()
    }
  }()

  icon, err := getIcon(fore)
  if err != nil {

    log.Fatal(err)
  }
  systray.SetIcon(icon)
  mUpdate := systray.AddMenuItem("Update Forecast", "Updates the forecast.io API")
  mSay := systray.AddMenuItem("Say Forecast", "Speak the forecast outloud")
  mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
  go func() {
    <-mQuit.ClickedCh
    systray.Quit()
    os.Exit(0)
  }()

  go func() {
    for {
      select {
      case <-mUpdate.ClickedCh:
        fore := getForecast(devkey, lat, long)
        icon, err := getIcon(fore)
        if err != nil {
          log.Fatal(err)
        }
        systray.SetIcon(icon)
        notifyForecast(fore)
      case <-mSay.ClickedCh:
        fore := getForecast(devkey, lat, long)
        msg := fmt.Sprintf("It feels like %.0f degrees.", fore.Currently.ApparentTemperature)
        msg2 := fmt.Sprintf("The weather is %s", fore.Currently.Summary)
        err = mack.Say(msg, "Samantha")
        err = mack.Say(msg2, "Samantha")
      }
    }
  }()

}

func getForecast(key string, lat string, long string) *forecast.Forecast {
  f, err := forecast.Get(key, lat, long, "now", forecast.US)
  if err != nil {
    log.Fatal(err)
  }
  return f
}

func getIcon(fore *forecast.Forecast) ([]byte, error) {
  var err error
  iconStr := strings.Replace(fore.Currently.Icon, "-", "", -1)
  iconAry := icon.IconMap[iconStr]
  if iconAry == nil {
    errMsg := fmt.Sprintf("No icon found for iconStr %s", iconStr)
    err = errors.New(errMsg)
  }
  return iconAry, err
}

func notifyForecast(fore *forecast.Forecast) {
  message := fmt.Sprintf("Feels like %.0fF", fore.Currently.ApparentTemperature)
  icon := fmt.Sprintf("icon/original/%s.png", fore.Currently.Icon)
  note := gosxnotifier.NewNotification(message)
  note.Title = "Forecast.io"
  note.Subtitle = fore.Currently.Summary
  note.AppIcon = icon
  note.ContentImage = icon
  note.Push()
}

func notifyRain() {
  note := gosxnotifier.NewNotification("Rain starting soon")
  note.Title = "Forecast.io"
  note.AppIcon = "icon/original/rain.png"
  note.ContentImage = "icon/original/rain.png"
  note.Push()
}

func checkForecast() {
  fore := getForecast(devkey, lat, long)
  if willRainSoon(fore) == true {
    icon, err := getIcon(fore)
    if err != nil {
      log.Fatal(err)
    }
    systray.SetIcon(icon)
    notifyRain()
  }
}

func willRainSoon(fore *forecast.Forecast) bool {
  tenMinutes := fore.Minutely.Data[0:10]
  for _, v := range tenMinutes {
    if v.PrecipProbability >= 0.7 {
      return true
    }
  }
  return false
}
