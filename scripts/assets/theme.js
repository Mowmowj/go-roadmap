// Theme toggle logic (loaded early to prevent FOUC)
(function () {
  var saved = localStorage.getItem("go-roadmap-theme");
  if (saved) {
    document.documentElement.setAttribute("data-theme", saved);
  }
})();

function initTheme() {
  var theme = localStorage.getItem("go-roadmap-theme") || "dark";
  document.documentElement.setAttribute("data-theme", theme);
  var btn = document.getElementById("themeToggle");
  if (btn) btn.textContent = theme === "dark" ? "☀️" : "🌙";

  // Update hljs theme link if present
  var hljsLink = document.getElementById("hljs-theme");
  if (hljsLink) {
    hljsLink.href =
      theme === "dark"
        ? "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css"
        : "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css";
  }
}

function toggleTheme() {
  var current = document.documentElement.getAttribute("data-theme") || "dark";
  var next = current === "dark" ? "light" : "dark";
  document.documentElement.setAttribute("data-theme", next);
  localStorage.setItem("go-roadmap-theme", next);
  var btn = document.getElementById("themeToggle");
  if (btn) btn.textContent = next === "dark" ? "☀️" : "🌙";

  var hljsLink = document.getElementById("hljs-theme");
  if (hljsLink) {
    hljsLink.href =
      next === "dark"
        ? "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css"
        : "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css";
  }
}
