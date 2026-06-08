document.addEventListener("DOMContentLoaded", async function() {
  var cabinet = document.querySelector(".cabinet");
  if (!cabinet) return;

  var searchForm = document.querySelector(".search-form");

  try {
    await API.me();
  } catch (error) {
    notify("Вы не авторизованы! Перенаправляем на страницу входа...", true);
    setTimeout(function() {
      window.location.href = "/login";
    }, 1000);
    return;
  }

  await loadLinks("");

  if (searchForm) {
    searchForm.addEventListener("submit", async function(e) {
      e.preventDefault();
      var query = searchForm.querySelector("input[name='q']").value.trim();
      await loadLinks(query);
    });
  }

  var logoutLink = document.querySelector("a[href='/login']");
  if (logoutLink) {
    logoutLink.addEventListener("click", async function(e) {
      e.preventDefault();
      await API.logout();
      notify("Вы вышли из системы");
      setTimeout(function() {
        window.location.href = "/login";
      }, 500);
    });
  }

  async function loadLinks(search) {
    if (search === undefined || search === null) search = "";
    try {
      var links = await API.getLinks(search);
      renderLinks(links);
    } catch (error) {
      notify("Ошибка загрузки: " + error, true);
    }
  }

  function renderLinks(links) {
    var oldRows = cabinet.querySelectorAll(".link-row");
    for (var i = 0; i < oldRows.length; i++) {
      oldRows[i].remove();
    }

    var oldEmpty = cabinet.querySelector(".empty-links-message");
    if (oldEmpty) oldEmpty.remove();

    if (links.length === 0) {
      var empty = document.createElement("p");
      empty.className = "empty-links-message";
      empty.style.cssText = "color:var(--muted); font-style:italic;";
      empty.textContent = "// ссылок пока нет. Создайте первую на главной!";
      cabinet.appendChild(empty);
      updateCounters(links);
      return;
    }

    for (var i = 0; i < links.length; i++) {
      var row = createLinkRow(links[i]);
      cabinet.appendChild(row);
    }

    updateCounters(links);
  }

  function createLinkRow(link) {
    var row = document.createElement("div");
    row.className = "link-row";
    row.setAttribute("data-id", link.id);

    var details = document.createElement("details");
    details.className = "link-details";

    var summary = document.createElement("summary");
    summary.className = "link-summary";

    var shortContainer = document.createElement("span");
    shortContainer.className = "link-short";
    shortContainer.style.display = "flex";
    shortContainer.style.alignItems = "center";
    shortContainer.style.gap = "8px";
    shortContainer.style.flexWrap = "wrap";

    var a = document.createElement("a");
    var fullShortUrl = window.location.protocol + "//" + window.location.host + link.short_url;
    a.href = fullShortUrl;
    a.textContent = link.short_url;
    a.style.textDecoration = "none";
    a.style.color = "var(--accent)";
    a.style.borderBottom = "2px solid var(--accent)";
    a.target = "_blank";
    
    var copyBtn = document.createElement("button");
    copyBtn.textContent = "Скопировать";
    copyBtn.style.cssText = 
      "background:none; border:1px solid var(--muted); cursor:pointer;" +
      "font-size:0.7rem; padding:2px 6px; border-radius:0;" +
      "font-family:var(--font); transition:0.1s;";
    copyBtn.title = "Копировать ссылку";
    copyBtn.addEventListener("click", function(e) {
      e.stopPropagation();
      navigator.clipboard.writeText(fullShortUrl).then(function() {
        var originalText = copyBtn.textContent;
        copyBtn.textContent = "✓";
        setTimeout(function() {
          copyBtn.textContent = originalText;
        }, 1500);
        notify("Ссылка скопирована!");
      }).catch(function() {
        notify("Не удалось скопировать", true);
      });
    });
    
    shortContainer.appendChild(a);
    shortContainer.appendChild(copyBtn);

    var long = document.createElement("span");
    long.className = "link-long";
    long.textContent = link.long_url;
    long.title = link.long_url;

    var clicks = document.createElement("span");
    clicks.className = "link-clicks";
    clicks.textContent = link.clicks;

    var date = document.createElement("span");
    date.className = "link-date";
    date.textContent = formatDate(link.created_at);

    var del = document.createElement("button");
    del.textContent = "✕";
    del.style.cssText =
      "background:none; border:2px solid var(--accent); cursor:pointer;" +
      "font-family:var(--font); font-weight:700; padding:2px 8px; color:var(--accent);";
    del.addEventListener("click", function(e) {
      e.stopPropagation();
      deleteLink(link.id, row);
    });

    summary.appendChild(shortContainer);
    summary.appendChild(long);
    summary.appendChild(clicks);
    summary.appendChild(date);
    summary.appendChild(del);
    details.appendChild(summary);

    var statsDiv = document.createElement("div");
    statsDiv.className = "link-stats";
    statsDiv.setAttribute("data-loaded", "false");

    details.addEventListener("toggle", async function() {
      if (details.open && statsDiv.getAttribute("data-loaded") === "false") {
        try {
          var stats = await API.getStats(link.id);
          renderStats(stats, statsDiv);
          statsDiv.setAttribute("data-loaded", "true");
        } catch (error) {
          statsDiv.innerHTML = "<p style='color:var(--muted);'>Статистика недоступна</p>";
        }
      }
    });

    details.appendChild(statsDiv);
    row.appendChild(details);
    return row;
  }

  function renderStats(stats, container) {
    if (!stats) {
      container.innerHTML = "<p style='color:var(--muted);'>Нет данных для отображения</p>";
      return;
    }
    var html = "";
    html += createStatColumn("Браузеры", stats.browsers);
    html += createStatColumn("Устройства", stats.devices);
    html += createStatColumn("Страны", stats.countries);
    html += createStatColumn("Источники", stats.referrers);
    container.innerHTML = html;
  }

  async function deleteLink(id, rowElement) {
    confirmDialog("Удалить ссылку?", async function() {
      try {
        await API.deleteLink(id);
        rowElement.remove();
        var links = await API.getLinks("");
        updateCounters(links);
        notify("Ссылка удалена");
      } catch (error) {
        notify("Ошибка: " + error, true);
      }
    });
  }

  function updateCounters(links) {
    var totalLinks = links.length;
    var totalClicks = 0;
    for (var i = 0; i < links.length; i++) {
      totalClicks += links[i].clicks;
    }

    var cards = document.querySelectorAll(".stat-card .num");
    if (cards.length >= 2) {
      cards[0].textContent = totalLinks;
      cards[1].textContent = totalClicks;
    }
  }
});