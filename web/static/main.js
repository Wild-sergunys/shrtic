document.addEventListener("DOMContentLoaded", function() {
  var form = document.querySelector(".shorten-form");
  if (!form) return;

  form.addEventListener("submit", async function(e) {
    e.preventDefault();

    var input = form.querySelector("input[name='url']");
    var url = input.value.trim();

    if (!url) {
      notify("Введите ссылку!", true);
      return;
    }

    try {
      var link = await API.createLink(url);
      showResult(link);
      input.value = "";
      notify("Ссылка создана!");
    } catch (error) {
      notify(error, true);
    }
  });

  function showResult(link) {
    var hero = document.querySelector(".hero");

    var oldResult = document.querySelector(".shorten-result");
    if (oldResult) oldResult.remove();

    var result = document.createElement("div");
    result.className = "shorten-result";
    result.style.cssText = "margin-top:24px; font-weight:700;";

    result.innerHTML =
      "<p style='color:var(--muted);'>// ссылка готова:</p>" +
      "<p style='font-size:1.2rem; margin-top:8px;'>" +
      "<a href='" + link.short_url + "' style='color:var(--accent);'>" +
      "localhost:8080" + link.short_url + "</a></p>";

    hero.appendChild(result);
  }
});