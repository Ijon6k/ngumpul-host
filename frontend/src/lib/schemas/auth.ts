import { z } from 'zod';

const passwordSchema = z.string().min(8, 'Password must be at least 8 characters.');

const usernameSchema = z
	.string()
	.min(1, 'Username is required.')
	.regex(
		/^[a-zA-Z0-9_.-]+$/,
		'Username may only contain letters, numbers, periods, underscores, and dashes.'
	);

const emailSchema = z.string().min(1, 'Email is required.').email('Enter a valid email address.');

export const loginSchema = z.object({
	emailOrUsername: z.string().min(1, 'Email or username is required.'),
	password: z.string().min(1, 'Password is required.')
});

export const registerSchema = z
	.object({
		invitationToken: z.string().optional().default(''),
		displayName: z.string().min(1, 'Display name is required.'),
		username: usernameSchema.max(32, 'Username must be 32 characters or fewer.'),
		email: emailSchema,
		password: passwordSchema,
		confirmPassword: z.string().min(1, 'Confirm password is required.')
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match.',
		path: ['confirmPassword']
	});

export const setupSchema = z
	.object({
		subdomain: z
			.string()
			.regex(/^[a-z0-9-]*$/, 'Subdomain may only contain lowercase letters, numbers, and hyphens.')
			.optional()
			.default(''),
		domain: z.string().min(1, 'Node domain is required.'),
		name: z
			.string()
			.min(1, 'Operator name is required.')
			.max(128, 'Operator name must be 128 characters or fewer.'),
		username: z
			.string()
			.min(3, 'Username must be at least 3 characters.')
			.max(30, 'Username must be 30 characters or fewer.')
			.regex(/^[a-zA-Z0-9_-]+$/, 'Username may only contain letters, numbers, underscores, and hyphens.'),
		email: emailSchema,
		password: passwordSchema,
		confirmPassword: z.string().min(1, 'Confirm password is required.')
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match.',
		path: ['confirmPassword']
	});

export type LoginData = z.infer<typeof loginSchema>;
export type RegisterData = z.infer<typeof registerSchema>;
export type SetupData = z.infer<typeof setupSchema>;