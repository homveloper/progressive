/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/**/*.templ',
    './internal/**/*.go',
  ],
  theme: {
    extend: {
      // 필요하면 여기에 커스텀 테마 추가 가능
    },
  },
  plugins: [
    // 필요한 플러그인 있으면 추가
  ],
}
