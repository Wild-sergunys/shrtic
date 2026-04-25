// Форматирование даты из ISO в ДД.ММ.ГГГГ
function formatDate(isoString) {
  var d = new Date(isoString);
  var day = String(d.getDate()).padStart(2, "0");
  var month = String(d.getMonth() + 1).padStart(2, "0");
  var year = d.getFullYear();
  return day + "." + month + "." + year;
}

// Кастомное уведомление
function notify(message, isError) {
  // Удаляем старое уведомление, если есть
  var old = document.querySelector(".shrtik-notify");
  if (old) old.remove();

  var div = document.createElement("div");
  div.className = "shrtik-notify";
  div.textContent = message;

  div.style.cssText =
    "position:fixed; top:20px; left:50%; transform:translateX(-50%);" +
    "z-index:9999; padding:12px 28px;" +
    "font-family:'Courier New',monospace; font-weight:700; font-size:0.9rem;" +
    "border:3px solid #1a1a1a; box-shadow:6px 6px 0 #1a1a1a;" +
    "background:" + (isError ? "#fff" : "#1a1a1a") + ";" +
    "color:" + (isError ? "#1a1a1a" : "#fff") + ";" +
    "max-width:90%; text-align:center; transition:opacity 0.3s;";

  document.body.appendChild(div);

  // Автоматически скрываем через 3 секунды
  setTimeout(function() {
    div.style.opacity = "0";
    setTimeout(function() {
      if (div.parentNode) div.remove();
    }, 300);
  }, 3000);
}

// Кастомный диалог подтверждения в стиле сайта
function confirmDialog(message, onYes, onNo) {
  // Удаляем старый диалог
  var old = document.querySelector(".shrtik-dialog");
  if (old) old.remove();

  var overlay = document.createElement("div");
  overlay.className = "shrtik-dialog";
  overlay.style.cssText =
    "position:fixed; top:0; left:0; width:100%; height:100%;" +
    "background:rgba(0,0,0,0.5); z-index:9998;" +
    "display:flex; align-items:center; justify-content:center;";

  var box = document.createElement("div");
  box.style.cssText =
    "background:#fff; border:3px solid #1a1a1a; box-shadow:6px 6px 0 #1a1a1a;" +
    "padding:24px 32px; max-width:400px; width:90%;" +
    "font-family:'Courier New',monospace; text-align:center;";

  var msg = document.createElement("p");
  msg.textContent = message;
  msg.style.cssText = "margin-bottom:20px; font-weight:700; font-size:0.95rem;";

  var btnYes = document.createElement("button");
  btnYes.textContent = "Да";
  btnYes.style.cssText =
    "padding:8px 24px; border:3px solid #1a1a1a; background:#1a1a1a; color:#fff;" +
    "font-family:'Courier New',monospace; font-weight:700; cursor:pointer;" +
    "margin-right:12px;";
  btnYes.addEventListener("click", function() {
    overlay.remove();
    if (onYes) onYes();
  });

  var btnNo = document.createElement("button");
  btnNo.textContent = "Нет";
  btnNo.style.cssText =
    "padding:8px 24px; border:3px solid #1a1a1a; background:#fff; color:#1a1a1a;" +
    "font-family:'Courier New',monospace; font-weight:700; cursor:pointer;";
  btnNo.addEventListener("click", function() {
    overlay.remove();
    if (onNo) onNo();
  });

  box.appendChild(msg);
  box.appendChild(btnYes);
  box.appendChild(btnNo);
  overlay.appendChild(box);
  document.body.appendChild(overlay);
}

// Создание колонки статистики (шахматная полоска + подписи)
function createStatColumn(title, items) {
  if (!items || items.length === 0) {
    return "<div class='stat-column'>" +
           "<h4>" + title + "</h4>" +
           "<p style='color:var(--muted);'>Нет данных</p>" +
           "</div>";
  }

  var html = "<div class='stat-column'>";
  html += "<h4>" + title + "</h4>";

  // Шахматная полоска
  html += "<div class='chess-bar'>";
  for (var i = 0; i < items.length; i++) {
    var colorClass = (i % 2 === 0) ? "chess-dark" : "chess-light";
    html += "<div class='chess-block " + colorClass + "' style='width:" + items[i].percentage + "%'></div>";
  }
  html += "</div>";

  // Подписи с процентами
  html += "<div class='perc-row'>";
  for (var j = 0; j < items.length; j++) {
    html += "<span>" + items[j].name + " " + items[j].percentage + "%</span>";
  }
  html += "</div>";

  html += "</div>";
  return html;
}