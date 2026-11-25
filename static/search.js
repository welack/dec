function escapeHTML(str) {
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

async function initSearch() {
  const input = document.getElementById("search-input");
  const results = document.getElementById("search-results");
  if (!input || !results) return;

  const params = new URLSearchParams(window.location.search);
  const initialQuery = (params.get("q") || "").trim();
  input.value = initialQuery;

  let posts = [];
  try {
    const res = await fetch("/search-index.json");
    if (!res.ok) throw new Error("Failed to load index");
    posts = await res.json();
  } catch (err) {
    results.innerHTML = "<p>Failed to load search index.</p>";
    return;
  }

  function render(query) {
    const term = query.trim().toLowerCase();
    if (!term) {
      results.innerHTML = "";
      return;
    }

    const matches = posts.filter((post) => {
      const title = (post.title || "").toLowerCase();
      const content = (post.content || "").toLowerCase();
      return title.includes(term) || content.includes(term);
    });

    if (matches.length === 0) {
      results.innerHTML = "<p>No results found.</p>";
      return;
    }

    results.innerHTML = matches
      .map((post) => `
        <article class="search-result">
          <time>${escapeHTML(post.date || "")}</time>
          <a href="${escapeHTML(post.url || "#")}">${escapeHTML(post.title || "Untitled")}</a>
        </article>
      `)
      .join("");
  }

  render(initialQuery);

  input.addEventListener("input", (e) => {
    const value = e.target.value;
    const url = new URL(window.location.href);

    if (value.trim()) url.searchParams.set("q", value.trim());
    else url.searchParams.delete("q");

    history.replaceState({}, "", url);
    render(value);
  });
}

document.addEventListener("DOMContentLoaded", initSearch);
