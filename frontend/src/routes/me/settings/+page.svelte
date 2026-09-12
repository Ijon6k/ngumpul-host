<script lang="ts">
	import { user } from '$lib/stores/auth';
	import { authApi, extractError } from '$lib/api';

	let displayName = $user?.display_name || '';
	let bio = $user?.bio || '';
	let avatarURL = $user?.avatar_url || '';
	let saving = false;
	let uploading = false;
	let error: string | null = null;
	let success: string | null = null;

	async function handleUpload(e: Event) {
		const target = e.target as HTMLInputElement;
		if (!target.files || target.files.length === 0) return;

		const file = target.files[0];
		uploading = true;
		error = null;

		try {
			const url = await authApi.uploadAvatar(file);
			if (url) {
				avatarURL = url;
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			uploading = false;
		}
	}

	async function handleSave() {
		saving = true;
		error = null;
		success = null;

		try {
			const res = await authApi.updateProfile({
				display_name: displayName,
				bio,
				avatar_url: avatarURL
			});

			if (res.user) {
				user.set(res.user);
			}
			success = 'Profile settings saved successfully.';
		} catch (err) {
			error = extractError(err);
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Profile Settings · Ngumpul Host</title>
</svelte:head>

<div class="p-6 sm:p-8 lg:p-10 max-w-3xl w-full flex flex-col gap-8">
	<header class="flex flex-col gap-1 pb-6 border-b border-(--border-hairline)">
		<h1 class="font-sans font-semibold text-2xl sm:text-3xl text-(--text-main) tracking-tight">Account & Profile Settings</h1>
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			Manage how your identity, avatar, and technical bio appear across the community roster.
		</p>
	</header>

	<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-lg p-6 sm:p-8 flex flex-col gap-6">
		{#if success}
			<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded text-xs font-mono">{success}</div>
		{/if}
		{#if error}
			<div class="p-3 bg-red-950/20 text-(--color-danger) border border-(--color-danger)/30 rounded text-xs font-mono">{error}</div>
		{/if}

		<form on:submit|preventDefault={handleSave} class="flex flex-col gap-5">
			<!-- Avatar Preview & Upload -->
			<div class="flex items-center gap-5 pb-5 border-b border-(--border-hairline)">
				{#if avatarURL}
					<img src={avatarURL} alt={displayName} class="w-14 h-14 rounded-full object-cover border border-(--border-hairline)" />
				{:else}
					<div class="w-14 h-14 rounded-full bg-(--accent-soft) text-(--accent-strong) text-lg font-mono font-bold flex items-center justify-center border border-(--border-hairline)">
						{(displayName || 'U').slice(0, 2).toUpperCase()}
					</div>
				{/if}

				<div class="flex flex-col gap-1">
					<label for="avatar-file" class="btn btn-secondary btn-sm cursor-pointer w-fit font-mono text-xs">
						{uploading ? 'Uploading...' : 'Upload Avatar'}
					</label>
					<input
						id="avatar-file"
						type="file"
						accept="image/jpeg,image/png,image/webp,image/gif"
						on:change={handleUpload}
						class="hidden"
						disabled={uploading}
					/>
					<span class="text-xs font-mono text-(--text-muted)">JPG, PNG, WebP up to 5MB.</span>
				</div>
			</div>

			<div class="flex flex-col gap-1">
				<label for="s-name" class="text-xs font-mono font-medium text-(--text-secondary)">Display Name</label>
				<input
					id="s-name"
					type="text"
					class="w-full px-3 py-2 bg-(--bg-muted)/40 border border-(--border-hairline) rounded text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
					required
					bind:value={displayName}
				/>
			</div>

			<div class="flex flex-col gap-1">
				<label for="s-bio" class="text-xs font-mono font-medium text-(--text-secondary)">Bio / Focus Area</label>
				<textarea
					id="s-bio"
					rows="3"
					class="w-full px-3 py-2 bg-(--bg-muted)/40 border border-(--border-hairline) rounded text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
					placeholder="Tell fellow engineers what you build or operate..."
					bind:value={bio}
				></textarea>
			</div>

			<button type="submit" class="btn btn-primary btn-sm self-start mt-1 font-mono text-xs" disabled={saving}>
				{saving ? 'Saving...' : 'Save Profile Settings'}
			</button>
		</form>
	</div>
</div>
