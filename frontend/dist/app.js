const openButton = document.getElementById("open-playlist");
const copyButton = document.getElementById("copy-titles");
const saveButton = document.getElementById("save-titles");
const saveStatus = document.getElementById("save-status");
const replaceChoice = document.getElementById("replace-choice");
const replaceButton = document.getElementById("replace-titles");
const cancelReplaceButton = document.getElementById("cancel-replace");
const status = document.getElementById("playlist-status");
const titleList = document.getElementById("titles");
const sanitizeButton = document.getElementById("sanitize-mp3");
const stopSanitizeButton = document.getElementById("stop-sanitize");
const sanitizeStatus = document.getElementById("sanitize-status");
const sanitizeErrors = document.getElementById("sanitize-errors");
let currentTitles = [];
let sanitizeRunning = false;
let stopRequested = false;

window.runtime.EventsOn("sanitize-progress", (message) => {
  if (!sanitizeRunning) return;
  stopSanitizeButton.hidden = false;
  sanitizeErrors.replaceChildren();
  if (!stopRequested) sanitizeStatus.textContent = message;
});

async function runSanitize() {
  sanitizeRunning = true;
  stopRequested = false;
  sanitizeButton.disabled = true;
  try {
    const result = await window.go.desktop.App.SanitizeMP3();
    if (result.cancelled) return;
    sanitizeErrors.replaceChildren();
    if (result.error) {
      sanitizeStatus.textContent = result.error;
      return;
    }
    const failures = result.failures ?? [];
    if (!result.stopped && result.updated + result.unchanged + failures.length === 0) {
      sanitizeStatus.textContent = "No MP3 files found.";
    } else {
      const prefix = result.stopped ? "Stopped · " : "";
      sanitizeStatus.textContent = `${prefix}Updated: ${result.updated} · Unchanged: ${result.unchanged} · Failed: ${failures.length}`;
    }
    for (const failure of failures.slice(0, 10)) {
      const item = document.createElement("li");
      item.textContent = `${failure.name}: ${failure.error}`;
      sanitizeErrors.append(item);
    }
    if (failures.length > 10) {
      const item = document.createElement("li");
      item.textContent = `…and ${failures.length - 10} more errors`;
      sanitizeErrors.append(item);
    }
  } catch (error) {
    sanitizeStatus.textContent = String(error);
  } finally {
    sanitizeRunning = false;
    sanitizeButton.disabled = false;
    stopSanitizeButton.hidden = true;
    stopSanitizeButton.disabled = false;
  }
}

sanitizeButton.addEventListener("click", runSanitize);
stopSanitizeButton.addEventListener("click", async () => {
  stopRequested = true;
  stopSanitizeButton.disabled = true;
  sanitizeStatus.textContent = "Stopping after current file…";
  try {
    await window.go.desktop.App.StopSanitize();
  } catch (error) {
    sanitizeStatus.textContent = String(error);
  }
});

copyButton.addEventListener("click", async () => {
  try {
    await window.go.desktop.App.CopyTitles(currentTitles);
  } catch (error) {
    status.textContent = String(error);
  }
});

async function saveTitles(replace) {
  saveButton.disabled = true;
  replaceButton.disabled = true;
  try {
    const result = await window.go.desktop.App.SaveTitles(replace);
    if (result.exists) {
      replaceChoice.hidden = false;
      return;
    }
    replaceChoice.hidden = true;
    saveStatus.textContent = result.error ? `Not saved: ${result.error}` : "Saved";
  } catch (error) {
    replaceChoice.hidden = true;
    saveStatus.textContent = `Not saved: ${String(error)}`;
  } finally {
    saveButton.disabled = false;
    replaceButton.disabled = false;
  }
}

saveButton.addEventListener("click", () => saveTitles(false));
replaceButton.addEventListener("click", () => saveTitles(true));
cancelReplaceButton.addEventListener("click", () => {
  replaceChoice.hidden = true;
  saveStatus.textContent = "Not saved";
});

openButton.addEventListener("click", async () => {
  openButton.disabled = true;
  try {
    const result = await window.go.desktop.App.OpenPlaylist();
    if (result.cancelled) return;

    titleList.replaceChildren();
    currentTitles = [];
    copyButton.disabled = true;
    saveButton.hidden = true;
    replaceChoice.hidden = true;
    saveStatus.textContent = "";
    if (result.error) {
      status.textContent = result.error;
      return;
    }

    currentTitles = result.titles ?? [];
    status.textContent = currentTitles.length === 0 ? "No titles found" : result.name;
    copyButton.disabled = currentTitles.length === 0;
    saveButton.hidden = false;
    saveStatus.textContent = "Not saved";
    for (const title of currentTitles) {
      const item = document.createElement("li");
      item.textContent = title;
      titleList.append(item);
    }
  } catch (error) {
    titleList.replaceChildren();
    currentTitles = [];
    copyButton.disabled = true;
    saveButton.hidden = true;
    replaceChoice.hidden = true;
    saveStatus.textContent = "";
    status.textContent = String(error);
  } finally {
    openButton.disabled = false;
  }
});
