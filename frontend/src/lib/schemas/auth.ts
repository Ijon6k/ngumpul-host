import * as v from 'valibot';

const password = v.pipe(
	v.string(),
	v.minLength(8, 'Password must be at least 8 characters.')
);

const username = v.pipe(
	v.string(),
	v.nonEmpty('Username is required.'),
	v.regex(
		/^[a-zA-Z0-9_.-]+$/,
		'Username may only contain letters, numbers, periods, underscores, and dashes.'
	)
);

const emailSchema = v.pipe(
	v.string(),
	v.nonEmpty('Email is required.'),
	v.email('Enter a valid email address.')
);

export const loginSchema = v.object({
	emailOrUsername: v.pipe(v.string(), v.nonEmpty('Email or username is required.')),
	password
});

export const registerSchema = v.object({
	invitationToken: v.optional(v.string(), ''),
	displayName: v.pipe(v.string(), v.nonEmpty('Display name is required.')),
	username: v.pipe(username, v.maxLength(32, 'Username must be 32 characters or fewer.')),
	email: emailSchema,
	password,
	confirmPassword: v.pipe(v.string(), v.minLength(8, 'Password must be at least 8 characters.'))
});

export const setupSchema = v.object({
	domain: v.pipe(v.string(), v.nonEmpty('Node domain is required.')),
	name: v.pipe(
		v.string(),
		v.nonEmpty('Operator name is required.'),
		v.maxLength(128, 'Operator name must be 128 characters or fewer.')
	),
	username: v.pipe(
		username,
		v.minLength(3, 'Username must be at least 3 characters.'),
		v.maxLength(30, 'Username must be 30 characters or fewer.'),
		v.regex(/^[a-zA-Z0-9_-]+$/, 'Username may only contain letters, numbers, underscores, and hyphens.')
	),
	email: emailSchema,
	password,
	confirmPassword: v.pipe(v.string(), v.minLength(8, 'Password must be at least 8 characters.'))
});

export type LoginData = v.InferOutput<typeof loginSchema>;
export type RegisterData = v.InferOutput<typeof registerSchema>;
export type SetupData = v.InferOutput<typeof setupSchema>;