# 2017 Solar Eclipse Traffic Visualization

See the finished product: https://youtu.be/AVwWofUFl3k

This project periodically fetches a map of the current traffic conditions
in the United States. I ran it during the day of the 2017 solar eclipse
to capture both the influx of traffic into the path of totality, along
with the outflux once the eclipse was over.

The screenshots are captured using the `phantomjs` headless web browser, with
the Selenium web driver interface to drive it. Some minor post-processing
to crop the images add timestamps is accomplished using the `draw2d` package.

To convert the video frames to a video, I used `avconv`:

```bash
$ avconv -r 12 -i 'frames/%d.png' -b:v 1000k eclipse_traffic.mp4
```

The music track is "Epic" from http://www.bensound.com/royalty-free-music.
I used `audacity` to cut it to the same length as the video and apply a
fade out effect at the end.

I used `avcon` once more to add the audio track:

```bash
$ avconv -strict -2 -i eclipse_traffic.mp4 -i bensound-epic-cut.mp3 eclipse_traffic_final.mp4
```
