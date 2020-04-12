package crawler

import "database/sql"

func ExecuteMigration2(db *sql.DB) error {
	var query = `
		ALTER TABLE profx.sources
			ADD COLUMN name VARCHAR(64) NOT NULL FIRST,
			ADD COLUMN kind VARCHAR(64) NOT NULL AFTER url;
		ALTER TABLE profx.sources
			DROP COLUMN rule;
		CREATE TABLE IF NOT EXISTS profx.source_rules
		(
			source  VARCHAR(64)                               NOT NULL,
			type    TINYINT UNSIGNED                          NOT NULL,
			rule    VARCHAR(1000)                             NOT NULL,
			created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
		);
	`
	_, err := db.Exec(query)

	return err
}

func RevertMigration2(db *sql.DB) error {
	var query = `
		ALTER TABLE profx.sources
			DROP COLUMN name,
			DROP COLUMN kind;
		ALTER TABLE profx.sources
			ADD COLUMN rule VARCHAR(1000) NOT NULL;
		DROP TABLE IF EXISTS profx.source_rules;
	`
	_, err := db.Exec(query)

	return err
}

func DataMigration2(db *sql.DB) error {
	var query1 = `INSERT INTO profx.sources (name, url, kind) VALUES (?, ?, ?)`
	var query2 = `INSERT INTO profx.source_rules (source, type, rule) VALUES (?, ?, ?)`
	data := []struct {
		Name string
		URL  string
		Rule string
	}{
		{
			Name: "mozilla", URL: "https://hacks.mozilla.org", Rule: `^https:\/{2}hacks\.mozilla\.org\/\d{4}\/\d{2}\/[a-zA-z0-9-]+\/$`,
		},
		{
			Name: "github", URL: "https://github.blog", Rule: `^https:\/{2}github\.blog\/\d{4}-\d{2}-\d{2}-[a-zA-z0-9-]+\/$`,
		},
		{
			Name: "dropbox", URL: "https://blogs.dropbox.com/tech/", Rule: `^https:\/{2}blogs\.dropbox\.com\/tech\/\d{4}\/\d{2}\/[a-zA-Z0-9-]+\/$`,
		},
		{
			Name: "spotify", URL: "https://labs.spotify.com/", Rule: `^https:\/{2}labs\.spotify\.com\/\d{4}\/\d{2}\/\d{2}\/[a-zA-Z0-9-]+\/$`,
		},
		{
			Name: "highscalability", URL: "http://highscalability.com/", Rule: `^http:\/{2}highscalability\.com\/blog\/\d{1,4}\/\d{1,2}\/\d{1,2}\/[a-zA-Z0-9-]+\.html$`,
		},
		{
			Name: "facebook", URL: "https://engineering.fb.com/", Rule: `^https:\/{2}engineering\.fb\.com\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/$`,
		},
		{
			Name: "linkedin", URL: "https://engineering.linkedin.com/", Rule: `^https:\/{2}engineering\.linkedin\.com\/blog\/\d{4}\/\d{2}\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "linkedin", URL: "https://engineering.linkedin.com/blog", Rule: `^https:\/{2}engineering\.linkedin\.com\/blog\/\d{4}\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "twilio", URL: "https://www.twilio.com/blog", Rule: `^https:\/{2}www\.twilio\.com\/blog\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "discord", URL: "https://blog.discordapp.com/tagged/engineering", Rule: `^https:\/{2}blog\.discordapp\.com\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "uber", URL: "https://eng.uber.com/", Rule: `^https:\/{2}eng\.uber\.com\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "booking", URL: "https://blog.booking.com/", Rule: `^https:\/{2}medium\.com\/booking-[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "booking", URL: "https://blog.booking.com/", Rule: `^https:\/{2}booking\.[a-zA-Z]+\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "netflixtechblog", URL: "https://netflixtechblog.com/", Rule: `^https:\/{2}netflixtechblog\.com\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "cloudera", URL: "https://blog.cloudera.com/category/technical/", Rule: `^https:\/{2}blog\.cloudera\.com\/[a-zA-Z0-9-]+\/$`,
		},
		{
			Name: "pinterest", URL: "https://medium.com/@Pinterest_Engineering", Rule: `^https:\/{2}medium\.com\/pinterest-engineering\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "buffer", URL: "https://open.buffer.com/", Rule: `^https:\/{2}open\.buffer\.com\/[a-zA-Z0-9-]+\/$`,
		},
		{
			Name: "deliveroo", URL: "https://deliveroo.engineering/articles/", Rule: `^https:\/{2}deliveroo\.engineering\/\d{4}\/\d{2}\/\d{2}\/[a-zA-Z0-9-]+\.html$`,
		},
		{
			Name: "mongo", URL: "https://www.mongodb.com/blog", Rule: `^https:\/{2}www\.mongodb\.com\/blog\/post\/[a-zA-Z0-9-]+$`,
		},
		{
			Name: "cloudflare", URL: "https://blog.cloudflare.com/", Rule: `^https:\/{2}blog\.cloudflare\.com\/[a-zA-Z0-9-]+\/$`,
		},
	}

	for _, item := range data {
		_, err := db.Exec(query1, item.Name, item.URL, "article")
		if err != nil {
			return err
		}

		_, err = db.Exec(query2, item.Name, 1, item.Rule)
		if err != nil {
			return err
		}
	}

	return nil
}
