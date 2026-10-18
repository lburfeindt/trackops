const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const vm = require("node:vm");

const page = fs.readFileSync(path.join(__dirname, "window-prototype.html"), "utf8");
const interactionModel = page.match(/<script>([\s\S]*?)<\/script>/)?.[1];
assert.ok(interactionModel, "prototype interaction model exists");
const { initialState, transition } = vm.runInNewContext(`${interactionModel}\n;({ initialState, transition })`);

test("canceling a second playlist choice restores the visible titles", () => {
  // Given
  const firstChoice = transition(initialState(), "OPEN_PLAYLIST");
  const playlistResult = transition(firstChoice, "PLAYLIST_TITLES");

  // When
  const secondChoice = transition(playlistResult, "OPEN_PLAYLIST");
  const cancelledChoice = transition(secondChoice, "CANCEL");

  // Then
  assert.equal(cancelledChoice.phase, "result");
  assert.equal(cancelledChoice.selection, "set.m3u8");
  assert.equal(cancelledChoice.titles.length, 3);
});

test("stopping ID3 processing preserves a partial result", () => {
  // Given
  const sourceChoice = transition(initialState(), "OPEN_ID3");
  const processing = transition(sourceChoice, "CHOOSE_FOLDER");

  // When
  const stopped = transition(processing, "ID3_CANCELLED");

  // Then
  assert.equal(processing.phase, "processing");
  assert.equal(stopped.phase, "result");
  assert.match(stopped.summary, /Earlier changes remain saved/);
});

test("empty and unreadable playlists can be retried", () => {
  for (const outcome of ["PLAYLIST_EMPTY", "PLAYLIST_ERROR"]) {
    // Given
    const choice = transition(initialState(), "OPEN_PLAYLIST");

    // When
    const result = transition(choice, outcome);
    const retry = transition(result, "OPEN_PLAYLIST");

    // Then
    assert.equal(retry.phase, "choose playlist");
    assert.equal(retry.tool, "Format Playlist");
  }
});

test("playlist selection exposes both titles and the sibling output file", () => {
  // Given
  const choice = transition(initialState(), "OPEN_PLAYLIST");

  // When
  const result = transition(choice, "PLAYLIST_TITLES");

  // Then
  assert.equal(result.phase, "result");
  assert.equal(result.titles.length, 3);
  assert.equal(result.outputFile, "set-titles.txt");
});

test("empty playlist writes an empty output file, while a save error keeps titles", () => {
  // Given
  const choice = transition(initialState(), "OPEN_PLAYLIST");

  // When
  const empty = transition(choice, "PLAYLIST_EMPTY");
  const failedSave = transition(choice, "PLAYLIST_WRITE_ERROR");

  // Then
  assert.equal(empty.outputFile, "empty-titles.txt");
  assert.equal(empty.titles.length, 0);
  assert.equal(failedSave.outputFile, "");
  assert.equal(failedSave.titles.length, 2);
  assert.match(failedSave.error, /Could not save/);
});
