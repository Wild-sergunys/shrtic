var MOCK = {
  // Авторизованный пользователь (null = не авторизован)
  user: null,

  // Все ссылки
  links: [
    {
      id: 1,
      short_url: "/x7kQ2",
      long_url: "https://github.com/username/awesome-project",
      clicks: 1247,
      created_at: "2024-11-11T12:00:00Z"
    },
    {
      id: 2,
      short_url: "/a9mN4",
      long_url: "https://stackoverflow.com/questions/12345",
      clicks: 856,
      created_at: "2026-04-19T15:30:00Z"
    },
    {
      id: 3,
      short_url: "/b3kP1",
      long_url: "https://golang.org/doc/tutorial/getting-started",
      clicks: 423,
      created_at: "2026-04-20T09:15:00Z"
    }
  ],

  // Статистика по ссылкам (ключ - id ссылки)
  stats: {
    1: {
      total_clicks: 1247,
      browsers: [
        { name: "Chrome", percentage: 62, count: 773 },
        { name: "Firefox", percentage: 18, count: 225 },
        { name: "Safari", percentage: 12, count: 150 },
        { name: "Other", percentage: 8, count: 99 }
      ],
      devices: [
        { name: "Desktop", percentage: 71, count: 885 },
        { name: "Mobile", percentage: 24, count: 299 },
        { name: "Tablet", percentage: 5, count: 63 }
      ],
      countries: [
        { name: "Россия", percentage: 45, count: 561 },
        { name: "США", percentage: 28, count: 349 },
        { name: "Германия", percentage: 15, count: 187 },
        { name: "Other", percentage: 12, count: 150 }
      ],
      referrers: [
        { name: "Прямой", percentage: 45, count: 561 },
        { name: "Twitter", percentage: 28, count: 349 },
        { name: "Telegram", percentage: 22, count: 274 },
        { name: "Other", percentage: 5, count: 63 }
      ]
    },
    2: {
      total_clicks: 856,
      browsers: [
        { name: "Chrome", percentage: 80, count: 685 },
        { name: "Firefox", percentage: 15, count: 128 },
        { name: "Other", percentage: 5, count: 43 }
      ],
      devices: [
        { name: "Desktop", percentage: 90, count: 770 },
        { name: "Mobile", percentage: 10, count: 86 }
      ],
      countries: [
        { name: "Россия", percentage: 60, count: 514 },
        { name: "США", percentage: 25, count: 214 },
        { name: "Other", percentage: 15, count: 128 }
      ],
      referrers: [
        { name: "Прямой", percentage: 70, count: 599 },
        { name: "Telegram", percentage: 30, count: 257 }
      ]
    },
    3: {
      total_clicks: 423,
      browsers: [
        { name: "Chrome", percentage: 55, count: 233 },
        { name: "Firefox", percentage: 30, count: 127 },
        { name: "Safari", percentage: 10, count: 42 },
        { name: "Other", percentage: 5, count: 21 }
      ],
      devices: [
        { name: "Desktop", percentage: 50, count: 212 },
        { name: "Mobile", percentage: 45, count: 190 },
        { name: "Tablet", percentage: 5, count: 21 }
      ],
      countries: [
        { name: "Россия", percentage: 70, count: 296 },
        { name: "Германия", percentage: 20, count: 85 },
        { name: "Other", percentage: 10, count: 42 }
      ],
      referrers: [
        { name: "Прямой", percentage: 40, count: 169 },
        { name: "Twitter", percentage: 35, count: 148 },
        { name: "Telegram", percentage: 25, count: 106 }
      ]
    }
  },

  // Счётчик для ID новых ссылок
  nextId: 4
};