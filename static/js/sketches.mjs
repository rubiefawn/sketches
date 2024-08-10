const audio = document.querySelector("audio");
const progress = document.getElementById("playback-progress");
const time = document.getElementById("playback-time");
const btn_main = document.getElementById("btn-main");
const name_main = document.getElementById("name-main");

function toggle_playback(override) {
  if (override ?? audio.paused()) audio.play(); else audio.pause();
}

function set_song(name, src) {
  audio.src = src;
  progress.max = audio.duration;
  name_main.innerHTML = "name";
}

audio.addEventListener("timeupdate", _ => {
  progress.value = audio.currentTime;
  const bm = `${Math.floor(audio.currentTime / 60)}`.padStart(2, '0');
  const bs = `${Math.floor(audio.currentTime % 60)}`.padStart(2, '0');
  const em = `${Math.floor(audio.duration / 60)}`.padStart(2, '0');
  const es = `${Math.floor(audio.duration % 60)}`.padStart(2, '0');
  time.innerHTML = `${bm}:${bs}/${em}:${es}`;
});

btn_main.addEventListener("click", _ => toggle_playback());
progress.addEventListener("input", _ => audio.fastSeek(progress.value));
document.querySelectorAll(".song").forEach(song => {
  const version_selector = song.querySelector(".song-version");
  const name = song.querySelector(".song-name").innerHTML;
  const btn = song.querySelector(".btn-play-song");
  btn.addEventListener("click", e => {
    if (audio.paused) {
      set_song(`${version_selector.selectedOptions[0].innerHTML} @ ${date}`, version_selector.value);
    }
    toggle_playback();
  });
});
