import {
	GithubLogo,
	GitlabLogo,
	GoogleDriveLogo,
	FigmaLogo,
	YoutubeLogo,
	TwitterLogo,
	DiscordLogo,
	Globe,
	FileText,
	ArrowUpRight,
	LinkSimple
} from 'phosphor-svelte';
export interface DetectedLinkInfo {
	icon: any;
	label: string;
	domain: string;
}

/**
 * Detects the service or domain type from a given URL to return an appropriate icon and label.
 */
export function detectLinkInfo(url: string, fallbackLabel = 'Link'): DetectedLinkInfo {
	if (!url) {
		return { icon: LinkSimple, label: fallbackLabel, domain: '' };
	}

	let domain = '';
	try {
		const parsed = new URL(url.startsWith('http://') || url.startsWith('https://') ? url : `https://${url}`);
		domain = parsed.hostname.toLowerCase().replace(/^www\./, '');
	} catch {
		domain = url.toLowerCase();
	}

	if (domain.includes('github.com')) {
		return { icon: GithubLogo, label: 'GitHub', domain };
	}
	if (domain.includes('gitlab.com')) {
		return { icon: GitlabLogo, label: 'GitLab', domain };
	}
	if (domain.includes('drive.google.com')) {
		return { icon: GoogleDriveLogo, label: 'Google Drive', domain };
	}
	if (domain.includes('docs.google.com')) {
		return { icon: FileText, label: 'Google Docs', domain };
	}
	if (domain.includes('figma.com')) {
		return { icon: FigmaLogo, label: 'Figma', domain };
	}
	if (domain.includes('youtube.com') || domain.includes('youtu.be')) {
		return { icon: YoutubeLogo, label: 'YouTube Demo', domain };
	}
	if (domain.includes('vimeo.com')) {
		return { icon: YoutubeLogo, label: 'Video Demo', domain };
	}
	if (domain.includes('twitter.com') || domain.includes('x.com')) {
		return { icon: TwitterLogo, label: 'X (Twitter)', domain };
	}
	if (domain.includes('discord.gg') || domain.includes('discord.com')) {
		return { icon: DiscordLogo, label: 'Discord', domain };
	}
	if (domain.includes('notion.so') || domain.includes('notion.site')) {
		return { icon: FileText, label: 'Notion Docs', domain };
	}

	return {
		icon: LinkSimple,
		label: fallbackLabel || domain || 'Website',
		domain
	};
}
