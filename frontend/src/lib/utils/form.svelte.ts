import type { z, ZodTypeAny } from 'zod';
import { extractError } from '$lib/api';

export interface FormOptions<TValues extends Record<string, any>> {
	schema: z.ZodType<TValues>;
	initialValues: TValues;
	onSubmit: (values: TValues) => Promise<void>;
}

export function createForm<TValues extends Record<string, any>>(options: FormOptions<TValues>) {
	let values = $state<TValues>({ ...options.initialValues });
	let errors = $state<Record<string, string>>({});
	let serverError = $state<string | null>(null);
	let isSubmitting = $state(false);

	function validate(): boolean {
		const result = options.schema.safeParse(values);
		if (!result.success) {
			const formatted: Record<string, string> = {};
			for (const issue of result.error.issues) {
				const field = issue.path[0];
				if (field !== undefined && typeof field === 'string' && !formatted[field]) {
					formatted[field] = issue.message;
				}
			}
			errors = formatted;
			return false;
		}
		errors = {};
		return true;
	}

	async function handleSubmit(e?: SubmitEvent) {
		if (e) {
			e.preventDefault();
		}

		serverError = null;

		const result = options.schema.safeParse(values);
		if (!result.success) {
			const formatted: Record<string, string> = {};
			for (const issue of result.error.issues) {
				const field = issue.path[0];
				if (field !== undefined && typeof field === 'string' && !formatted[field]) {
					formatted[field] = issue.message;
				}
			}
			errors = formatted;
			return;
		}

		errors = {};
		isSubmitting = true;

		try {
			await options.onSubmit(result.data);
		} catch (err) {
			serverError = extractError(err);
		} finally {
			isSubmitting = false;
		}
	}

	function reset() {
		values = { ...options.initialValues };
		errors = {};
		serverError = null;
		isSubmitting = false;
	}

	return {
		get values() {
			return values;
		},
		set values(v: TValues) {
			values = v;
		},
		get errors() {
			return errors;
		},
		set errors(e: Record<string, string>) {
			errors = e;
		},
		get serverError() {
			return serverError;
		},
		set serverError(s: string | null) {
			serverError = s;
		},
		get isSubmitting() {
			return isSubmitting;
		},
		validate,
		handleSubmit,
		reset
	};
}
