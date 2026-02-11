import createMiddleware from 'next-intl/middleware';
import { locales } from './i18n';

export default createMiddleware({
  locales,
  defaultLocale: 'en',
  localePrefix: 'as-needed',
});

export const config = {
  matcher: ['/', '/(en|ms|zh-CN)/:path*', '/((?!api|_next|_vercel|.*\\..*).*)'],
};
