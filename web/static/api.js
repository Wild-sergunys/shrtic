var API = {
  _baseURL: '',

  _request: async function(method, path, body) {
    var headers = { 'Content-Type': 'application/json' };
    var token = window.localStorage.getItem('shrtic_token');
    if (token) {
      headers['Authorization'] = 'Bearer ' + token;
    }

    var options = {
      method: method,
      headers: headers
    };
    if (body) {
      options.body = JSON.stringify(body);
    }

    var response = await fetch(this._baseURL + path, options);
    var data = await response.json();

    if (!response.ok) {
      throw data.message || 'Ошибка сервера';
    }

    return data;
  },

  register: function(login, password) {
    return this._request('POST', '/api/auth/register', { login: login, password: password });
  },

  login: function(login, password) {
    return this._request('POST', '/api/auth/login', { login: login, password: password });
  },

  logout: function() {
    return this._request('POST', '/api/auth/logout');
  },

  me: function() {
    return this._request('GET', '/api/auth/me');
  },

  createLink: function(url) {
    return this._request('POST', '/api/links', { url: url });
  },

  getLinks: function(search) {
    var path = '/api/links';
    if (search) {
      path += '?search=' + encodeURIComponent(search);
    }
    return this._request('GET', path);
  },

  deleteLink: function(id) {
    return this._request('DELETE', '/api/links/' + id);
  },

  getStats: function(id) {
    return this._request('GET', '/api/links/' + id + '/stats');
  }
};