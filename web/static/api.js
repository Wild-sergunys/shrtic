var API = {
  // Базовая задержка для имитации сети
  _delay: 300,

  // Вспомогательная функция - имитация ответа сервера
  _respond: function(data, shouldReject) {
    var self = this;
    return new Promise(function(resolve, reject) {
      setTimeout(function() {
        if (shouldReject) {
          reject(data); // data = сообщение об ошибке
        } else {
          resolve(data); // data = полезные данные
        }
      }, self._delay);
    });
  },

  // Регистрация
  register: function(login, password) {
    if (login === "admin" || login === "user") {
      return this._respond("Пользователь с таким логином уже существует", true);
    }
    if (password.length < 6) {
      return this._respond("Пароль должен быть не менее 6 символов", true);
    }
    MOCK.user = { id: 1, login: login, created_at: new Date().toISOString() };
    return this._respond({ id: 1 });
  },

  // Вход
  login: function(login, password) {
    if (!login || !password) {
      return this._respond("Логин и пароль обязательны", true);
    }
    if (password.length < 6) {
      return this._respond("Неверный логин или пароль", true);
    }
    MOCK.user = { id: 1, login: login, created_at: new Date().toISOString() };
    return this._respond({
      token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock-token",
      role: "user"
    });
  },

  // Выход
  logout: function() {
    MOCK.user = null;
    return this._respond({ message: "Выход выполнен успешно" });
  },

  // Создать короткую ссылку
  createLink: function(url) {
    if (!url || url.length < 5) {
      return this._respond("Некорректный URL", true);
    }

    // Генерируем случайный код из 5 символов
    var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    var code = "";
    for (var i = 0; i < 5; i++) {
      code += chars.charAt(Math.floor(Math.random() * chars.length));
    }

    var link = {
      id: MOCK.nextId++,
      short_url: "/" + code,
      long_url: url,
      clicks: 0,
      created_at: new Date().toISOString()
    };

    MOCK.links.unshift(link);

    // Пустая статистика для новой ссылки
    MOCK.stats[link.id] = {
      total_clicks: 0,
      browsers: [],
      devices: [],
      countries: [],
      referrers: []
    };

    return this._respond(link);
  },

  // Получить список ссылок (с поиском)
  getLinks: function(search) {
    var result = MOCK.links;
    if (search) {
      var q = search.toLowerCase();
      result = result.filter(function(link) {
        return link.long_url.toLowerCase().indexOf(q) !== -1;
      });
    }
    return this._respond(result);
  },

  // Удалить ссылку
  deleteLink: function(id) {
    var found = false;
    for (var i = 0; i < MOCK.links.length; i++) {
      if (MOCK.links[i].id === id) {
        MOCK.links.splice(i, 1);
        delete MOCK.stats[id];
        found = true;
        break;
      }
    }
    if (found) {
      return this._respond({ message: "Ссылка удалена" });
    }
    return this._respond("Ссылка не найдена", true);
  },

  // Получить статистику по ссылке
  getStats: function(id) {
    var stats = MOCK.stats[id];
    if (stats) {
      return this._respond(stats);
    }
    return this._respond("Статистика не найдена", true);
  }
};