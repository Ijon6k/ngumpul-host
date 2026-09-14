import type { Cookies } from '@sveltejs/kit';

const SESSION_COOKIE_NAME = 'ngumpul_session';

export function forwardSessionCookie(cookies: Cookies, response: Response) {
	const setCookieHeaders = response.headers.getSetCookie();

	for (const raw of setCookieHeaders) {
		const segment = raw.split(';')[0]?.trim();
		const name = segment?.split('=')[0]?.trim();
		if (name !== SESSION_COOKIE_NAME || !segment) continue;

		const parts = segment.split('=');
		const value = parts.slice(1).join('=');
		if (!value) continue;

		const attrs = raw.slice(segment.length);
		const expiresMatch = attrs.match(/expires=([^;]+)/i);
		const maxAgeMatch = attrs.match(/max-age=([^;]+)/i);
		const secure = /\bsecure\b/i.test(attrs);
		const sameSiteMatch = attrs.match(/samesite=([^;]+)/i);

		const sameSite =
			sameSiteMatch?.[1]?.trim().toLowerCase() === 'strict' ? 'strict' : 'lax';

		cookies.set(SESSION_COOKIE_NAME, value, {
			path: '/',
			secure,
			sameSite,
			httpOnly: true,
			maxAge: maxAgeMatch ? parseInt(maxAgeMatch[1], 10) : undefined,
			expires: expiresMatch ? new Date(expiresMatch[1]) : undefined
		});
		break;
	}
}
