package main 

import (
	"github.com/elgopher/pi/piaudio"
	_ "embed"
)

//go:embed "assets/sounds/bat_hurt.wav"
var batHurtWAV []byte
var BatHurtSample *piaudio.Sample

//go:embed "assets/sounds/collect.wav"
var collectWAV []byte
var CollectSample *piaudio.Sample

//go:embed "assets/sounds/crab_hurt.wav"
var crabHurtWAV []byte
var CrabHurtSample *piaudio.Sample

//go:embed "assets/sounds/crab_walk.wav"
var crabWalkWAV []byte
var CrabWalkSample *piaudio.Sample

//go:embed "assets/sounds/explode.wav"
var explodeWAV []byte
var ExplodeSample *piaudio.Sample

//go:embed "assets/sounds/game_start.wav"
var gameStartWAV []byte
var GameStartSample *piaudio.Sample

//go:embed "assets/sounds/player_hurt.wav"
var playerHurtWAV []byte
var PlayerHurtSample *piaudio.Sample

//go:embed "assets/sounds/player_walk.wav"
var playerWalkWAV []byte
var PlayerWalkSample *piaudio.Sample

//go:embed "assets/sounds/projectile.wav"
var projectileWAV []byte
var ProjectileSample *piaudio.Sample

//go:embed "assets/sounds/td_song_1.wav"
var themeSongWAV []byte 
var ThemeSongSample *piaudio.Sample 

func LoadSamples() {
	BatHurtSample = piaudio.DecodeWav(batHurtWAV)
	CollectSample = piaudio.DecodeWav(collectWAV)
	CrabHurtSample = piaudio.DecodeWav(crabHurtWAV)
	CrabWalkSample = piaudio.DecodeWav(crabWalkWAV)
	ExplodeSample = piaudio.DecodeWav(explodeWAV)
	GameStartSample = piaudio.DecodeWav(gameStartWAV)
	PlayerHurtSample = piaudio.DecodeWav(playerHurtWAV)
	PlayerWalkSample = piaudio.DecodeWav(playerWalkWAV)
	ProjectileSample = piaudio.DecodeWav(projectileWAV)
	ThemeSongSample = piaudio.DecodeWav(themeSongWAV)

	piaudio.LoadSample(BatHurtSample)
	piaudio.LoadSample(CollectSample)
	piaudio.LoadSample(CrabHurtSample)
	piaudio.LoadSample(CrabWalkSample)
	piaudio.LoadSample(ExplodeSample)
	piaudio.LoadSample(GameStartSample)
	piaudio.LoadSample(PlayerHurtSample)
	piaudio.LoadSample(PlayerWalkSample)
	piaudio.LoadSample(ProjectileSample)
	piaudio.LoadSample(ThemeSongSample)	
}

func PlaySound(sample *piaudio.Sample) {
	ch := piaudio.Chan1 | piaudio.Chan2
	vol := 1.0
	pitch := 1.0 
	piaudio.Play(ch, sample, pitch, vol)
}

func PlayTheme() {
	ch := piaudio.Chan3 | piaudio.Chan4
	vol := .20
	pitch := 1.0 
	piaudio.Play(ch, ThemeSongSample, pitch, vol)	
	piaudio.SetLoop(ch, 0, ThemeSongSample.Len(), piaudio.LoopForward, 0)
}

func StopTheme() {
	ch := piaudio.Chan3 | piaudio.Chan4
	vol := 0.0
	pitch := 1.0 
	piaudio.Play(ch, ThemeSongSample, pitch, vol)	
}