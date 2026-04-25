document.addEventListener("DOMContentLoaded", async function() {
  var cabinet = document.querySelector(".cabinet");
  if (!cabinet) return;

  // Проверка авторизации
  var token = localStorage.getItem("shrtik_token");
  if (!token) {
    notify("Вы не авторизованы! Перенаправляем на страницу входа...", true);
    setTimeout(function() {
      window.location.href = "auth.html";
    }, 1000);
    return;
  }

  var searchForm = document.querySelector(".search-form");

  // Загружаем ссылки при открытии
  await loadLinks("");

  // Поиск
  if (searchForm) {
    searchForm.addEventListener("submit", async function(e) {
      e.preventDefault();
      var query = searchForm.querySelector("input[name='q']").value.trim();
      await loadLinks(query);
    });
  }

  // Выход
  var logoutLink = document.querySelector("a[href='auth.html']");
  if (logoutLink) {
    logoutLink.addEventListener("click", function(e) {
      e.preventDefault();
      localStorage.removeItem("shrtik_token");
      localStorage.removeItem("shrtik_login");
      notify("Вы вышли из системы");
      setTimeout(function() {
        window.location.href = "auth.html";
      }, 500);
    });
  }

  // Загрузка ссылок
  async function loadLinks(search) {
    try {
      var links = await API.getLinks(search);
      renderLinks(links);
    } catch (error) {
      notify("Ошибка загрузки: " + error, true);
    }
  }

  // Отрисовка списка
  function renderLinks(links) {
    // Удаляем старые строки
    var oldRows = cabinet.querySelectorAll(".link-row");
    for (var i = 0; i < oldRows.length; i++) {
      oldRows[i].remove();
    }

    // Пустой список
    if (links.length === 0) {
      var empty = document.createElement("p");
      empty.style.cssText = "color:var(--muted); font-style:italic;";
      empty.textContent = "// ссылок пока нет. Создайте первую на главной!";
      cabinet.appendChild(empty);
      updateCounters(links);
      return;
    }

    // Создаём строки
    for (var i = 0; i < links.length; i++) {
      var row = createLinkRow(links[i]);
      cabinet.appendChild(row);
    }

    updateCounters(links);
  }

  // Создание строки ссылки
  function createLinkRow(link) {
    var row = document.createElement("div");
    row.className = "link-row";
    row.setAttribute("data-id", link.id);

    var details = document.createElement("details");
    details.className = "link-details";

    // Заголовок
    var summary = document.createElement("summary");
    summary.className = "link-summary";

    // Короткая ссылка
    var short = document.createElement("span");
    short.className = "link-short";
    var a = document.createElement("a");
    a.href = link.short_url;
    a.textContent = link.short_url;
    short.appendChild(a);

    // Длинный URL
    var long = document.createElement("span");
    long.className = "link-long";
    long.textContent = link.long_url;

    // Клики
    var clicks = document.createElement("span");
    clicks.className = "link-clicks";
    clicks.textContent = link.clicks;

    // Дата
    var date = document.createElement("span");
    date.className = "link-date";
    date.textContent = formatDate(link.created_at);

    // Кнопка удаления
    var del = document.createElement("button");
    del.textContent = "✕";
    del.style.cssText =
      "background:none; border:2px solid var(--accent); cursor:pointer;" +
      "font-family:var(--font); font-weight:700; padding:2px 8px; color:var(--accent);";
    del.addEventListener("click", function(e) {
      e.stopPropagation();
      deleteLink(link.id, row);
    });

    summary.appendChild(short);
    summary.appendChild(long);
    summary.appendChild(clicks);
    summary.appendChild(date);
    summary.appendChild(del);
    details.appendChild(summary);

    // Статистика (ленивая загрузка)
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

  // Отрисовка статистики
  function renderStats(stats, container) {
    var html = "";
    html += createStatColumn("Браузеры", stats.browsers);
    html += createStatColumn("Устройства", stats.devices);
    html += createStatColumn("Страны", stats.countries);
    html += createStatColumn("Источники", stats.referrers);
    container.innerHTML = html;
  }

  // Удаление ссылки
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

  // Обновление счётчиков
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