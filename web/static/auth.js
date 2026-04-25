document.addEventListener("DOMContentLoaded", function() {
  var container = document.querySelector(".auth-container");
  if (!container) return;

  var loginForm = document.querySelector(".tab-login form");
  var registerForm = document.querySelector(".tab-register form");

  // Если уже авторизован — предлагаем перейти в кабинет
  var token = localStorage.getItem("shrtik_token");
  if (token) {
    confirmDialog("Вы уже авторизованы. Перейти в кабинет?", function() {
      window.location.href = "cabinet.html";
    });
  }

  // Обработчик входа
  if (loginForm) {
    loginForm.addEventListener("submit", async function(e) {
      e.preventDefault();

      var login = loginForm.querySelector("input[name='login']").value.trim();
      var password = loginForm.querySelector("input[name='password']").value.trim();

      if (!login || !password) {
        notify("Заполните все поля!", true);
        return;
      }

      try {
        var response = await API.login(login, password);
        localStorage.setItem("shrtik_token", response.token);
        localStorage.setItem("shrtik_login", login);
        notify("Вход выполнен! Переходим в кабинет...");
        setTimeout(function() {
          window.location.href = "cabinet.html";
        }, 500);
      } catch (error) {
        notify(error, true);
      }
    });
  }

  // Обработчик регистрации
  if (registerForm) {
    registerForm.addEventListener("submit", async function(e) {
      e.preventDefault();

      var login = registerForm.querySelector("input[name='login']").value.trim();
      var password = registerForm.querySelector("input[name='password']").value.trim();

      if (!login || !password) {
        notify("Заполните все поля!", true);
        return;
      }

      try {
        await API.register(login, password);
        notify("Регистрация успешна! Теперь войдите.");

        document.getElementById("tab-login").checked = true;
        loginForm.querySelector("input[name='login']").value = login;
        registerForm.reset();
      } catch (error) {
        notify(error, true);
      }
    });
  }
});