import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'updock',
  description: 'Run any Docker app from one word.',
  lang: 'en-US',
  base: '/updock/',
  srcDir: 'src',
  cleanUrls: true,
  lastUpdated: true,
  appearance: 'dark',
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/updock/logo.svg' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Hanken+Grotesk:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600;700&display=swap' }],
    ['link', { rel: 'stylesheet', href: 'https://cdn.jsdelivr.net/npm/@fortawesome/fontawesome-free@6.7.2/css/all.min.css' }],
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Installation', link: '/installation' },
      { text: 'Quick start', link: '/quickstart' },
      { text: 'Commands', link: '/commands' },
    ],
    sidebar: [
      {
        text: 'Getting started',
        items: [
          { text: 'Installation', link: '/installation' },
          { text: 'Quick start', link: '/quickstart' },
        ],
      },
      {
        text: 'Guide',
        items: [
          { text: 'Usage', link: '/usage' },
          { text: 'Commands', link: '/commands' },
          { text: 'Configuration', link: '/configuration' },
          { text: 'Recipes', link: '/recipes' },
          { text: 'Projects and file layout', link: '/projects' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'Security model', link: '/security' },
          { text: 'Troubleshooting', link: '/troubleshooting' },
          { text: 'FAQ', link: '/faq' },
        ],
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/amrelsagaei/updock' },
    ],
    search: {
      provider: 'local',
    },
    editLink: {
      pattern: 'https://github.com/amrelsagaei/updock/edit/main/docs/src/:path',
      text: 'Edit this page on GitHub',
    },
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright (c) 2026 Amr Elsagaei',
    },
  },
})
