package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
)

func SeedInitialData(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config) error {
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.Println("Seeding initial database data...")

	// 1. Seed Admin User
	adminPassHash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}

	var adminID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, role, status, bio)
		VALUES ($1, $2, $3, $4, 'ADMIN', 'ACTIVE', 'Host maintainer and community organizer.')
		RETURNING id
	`, cfg.AdminUsername, cfg.AdminEmail, adminPassHash, cfg.AdminName).Scan(&adminID)
	if err != nil {
		return err
	}

	// 2. Seed Community Member: Erik
	erikPass, _ := auth.HashPassword("erik123456")
	var erikID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, role, status, bio)
		VALUES ('erik', 'erik@example.com', $1, 'Erik', 'USER', 'ACTIVE', 'Building things and breaking them.')
		RETURNING id
	`, erikPass).Scan(&erikID)
	if err != nil {
		return err
	}

	// 3. Seed Community Member: Raka
	rakaPass, _ := auth.HashPassword("raka123456")
	var rakaID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, role, status, bio)
		VALUES ('raka', 'raka@example.com', $1, 'Raka', 'USER', 'ACTIVE', 'Designing interfaces, writing minimal backend services.')
		RETURNING id
	`, rakaPass).Scan(&rakaID)
	if err != nil {
		return err
	}

	// 4. Seed Community Member: Nanda
	nandaPass, _ := auth.HashPassword("nanda123456")
	var nandaID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, role, status, bio)
		VALUES ('nanda', 'nanda@example.com', $1, 'Nanda', 'USER', 'ACTIVE', 'Exploring distributed systems, microcontrollers, and Go.')
		RETURNING id
	`, nandaPass).Scan(&nandaID)
	if err != nil {
		return err
	}

	// 5. Seed Projects (Atlas, Discord Bot, Raka's Portfolio)
	var atlasID string
	err = pool.QueryRow(ctx, `
		INSERT INTO projects (
			owner_id, name, slug, description, cover_image_url, repository_url, documentation_url,
			demo_url, technology_stack, hosting_type, public_url, status, visibility, published_at
		) VALUES (
			$1, 'Atlas', 'atlas', 'A small collaborative web application for team note-taking and coordination.',
			'/assets/project_atlas_cover.jpg',
			'https://github.com/example/atlas', 'https://docs.example.com/atlas',
			'https://atlas.example.com', ARRAY['Svelte', 'Go', 'PostgreSQL'], 'HOSTED_HERE',
			'https://atlas.example.com', 'ONLINE', 'PUBLIC', NOW()
		) RETURNING id
	`, erikID).Scan(&atlasID)
	if err != nil {
		return err
	}

	var botID string
	err = pool.QueryRow(ctx, `
		INSERT INTO projects (
			owner_id, name, slug, description, cover_image_url, repository_url,
			technology_stack, hosting_type, public_url, status, visibility, published_at
		) VALUES (
			$1, 'Discord Bot', 'discord-bot', 'Lightweight community utility bot automating role assignments and notifications.',
			'/assets/project_bot_cover.jpg',
			'https://github.com/example/discord-bot', ARRAY['Go', 'Docker'], 'HOSTED_HERE',
			'https://bot.example.com', 'ONLINE', 'PUBLIC', NOW()
		) RETURNING id
	`, erikID).Scan(&botID)
	if err != nil {
		return err
	}

	var portfolioID string
	err = pool.QueryRow(ctx, `
		INSERT INTO projects (
			owner_id, name, slug, description, cover_image_url, repository_url,
			technology_stack, hosting_type, public_url, status, visibility, published_at
		) VALUES (
			$1, 'Raka Portfolio', 'raka-portfolio', 'Editorial personal portfolio showcasing design works and typography experiments.',
			'/assets/infra_abstract_cover.jpg',
			'https://github.com/raka/portfolio', ARRAY['SvelteKit', 'CSS'], 'EXTERNAL',
			'https://raka.me', 'ONLINE', 'PUBLIC', NOW()
		) RETURNING id
	`, rakaID).Scan(&portfolioID)
	if err != nil {
		return err
	}

	// 6. Seed System Status
	_, err = pool.Exec(ctx, `
		INSERT INTO system_status (id, name, status, response_time_ms) VALUES
		('portal', 'Portal Frontend', 'OPERATIONAL', 24),
		('api', 'REST API', 'OPERATIONAL', 18),
		('projects', 'Hosted Projects', 'OPERATIONAL', 45)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return err
	}

	// 7. Seed Activities (Safe, curated application events as defined in overview.md)
	_, err = pool.Exec(ctx, `
		INSERT INTO activities (actor_id, project_id, type, metadata, visibility, created_at) VALUES
		($1, $2, 'PROJECT_DEPLOYED', '{"title": "Erik deployed Atlas", "message": "Updated frontend components"}'::jsonb, 'PUBLIC', NOW() - INTERVAL '40 minutes'),
		($3, $4, 'PROJECT_PUBLISHED', '{"title": "Raka published a project", "message": "Raka Portfolio is now online"}'::jsonb, 'PUBLIC', NOW() - INTERVAL '3 hours'),
		($5, NULL, 'MEMBER_JOINED', '{"title": "Nanda joined Ngumpul", "message": "Welcome Nanda to the server"}'::jsonb, 'PUBLIC', NOW() - INTERVAL '7 hours'),
		($1, $6, 'PROJECT_UPDATED', '{"title": "Discord Bot was updated", "message": "Added role sync command"}'::jsonb, 'PUBLIC', NOW() - INTERVAL '1 day')
	`, erikID, atlasID, rakaID, portfolioID, nandaID, botID)
	if err != nil {
		return err
	}

	// 8. Seed Sample Notification for Erik
	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, title, body, data) VALUES
		($1, 'PROJECT_STATUS_CHANGED', 'Your project Atlas is online', 'Atlas is operational at https://atlas.example.com', '{"project_slug": "atlas"}'::jsonb)
	`, erikID)
	if err != nil {
		return err
	}

	log.Println("Initial data seeding completed successfully.")
	return nil
}
