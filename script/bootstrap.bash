psql -U postgres -d sendin <<'SQL'

-- ACCOUNT

INSERT INTO accounts (
    name,
    email,
    profile,
    picture,
    timezone,
    token,
    user_agent
)
VALUES (
    'Vishnu Shon',
    'vishnu.shon@example.com',
    'https://linkedin.com/in/vishnushon',
    'https://randomuser.me/api/portraits/men/4.jpg',
    'Asia/Kolkata',
    'AQEDARfake123vishnuTOKENxyz',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X)'
);

-- CONNECTIONS

INSERT INTO connections (
    account_id,
    first_name,
    last_name,
    public_id,
    bio,
    picture,
    company,
    country
)
SELECT
    id,
    'Arjun',
    'Sharma',
    'arjunsharma',
    'Senior Software Engineer building distributed systems and search products.',
    'https://randomuser.me/api/portraits/men/1.jpg',
    'Flipkart',
    'India'
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO connections (
    account_id,
    first_name,
    last_name,
    public_id,
    bio,
    picture,
    company,
    country
)
SELECT
    id,
    'Rohan',
    'Verma',
    'rohanverma',
    'Product Manager focused on growth, onboarding, and experimentation.',
    'https://randomuser.me/api/portraits/men/2.jpg',
    'Razorpay',
    'India'
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO connections (
    account_id,
    first_name,
    last_name,
    public_id,
    bio,
    picture,
    company,
    country
)
SELECT
    id,
    'Vikram',
    'Kapoor',
    'vikramkapoor',
    'Machine Learning Engineer working on recommendation systems.',
    'https://randomuser.me/api/portraits/men/3.jpg',
    'Adobe',
    'Singapore'
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO connections (
    account_id,
    first_name,
    last_name,
    public_id,
    bio,
    picture,
    company,
    country
)
SELECT
    id,
    'Amit',
    'Mehta',
    'amitmehta',
    'Product Designer passionate about developer tooling.',
    'https://randomuser.me/api/portraits/men/4.jpg',
    'Atlassian',
    'Australia'
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO connections (
    account_id,
    first_name,
    last_name,
    public_id,
    bio,
    picture,
    company,
    country
)
SELECT
    id,
    'Karan',
    'Gupta',
    'karangupta',
    'Marketing and growth lead scaling B2B SaaS products.',
    'https://randomuser.me/api/portraits/men/5.jpg',
    'HubSpot',
    'United States'
FROM accounts
WHERE email = 'vishnu.shon@example.com';

-- TEMPLATES

INSERT INTO templates (
    label,
    value
)
VALUES
(
    'Job Opportunity',
    'Hi {{name}}, I came across your profile and thought your background was a strong fit for a role we are hiring for.'
),
(
    'Follow Up',
    'Hey {{name}}, just following up on my previous message.'
),
(
    'Networking',
    'Hi {{name}}, I enjoyed reading about your work at {{company}} and would love to connect.'
),
(
    'Collaboration',
    'Hi {{name}}, I think there may be an opportunity for us to collaborate.'
),
(
    'Founder Outreach',
    'Hi {{name}}, I am building SendIn and wanted to reach out personally.'
);

-- MESSAGES

INSERT INTO messages (
    account_id,
    name,
    picture,
    profile,
    company,
    timezone,
    message,
    is_sent
)
SELECT
    id,
    'Sophia Johnson',
    'https://randomuser.me/api/portraits/women/21.jpg',
    'https://linkedin.com/in/sophiajohnson',
    'Predictable Hiring',
    'America/New_York',
    'Just checking in to see if you had a chance to review my last message.',
    false
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO messages (
    account_id,
    name,
    picture,
    profile,
    company,
    timezone,
    message,
    is_sent
)
SELECT
    id,
    'Olivia Brown',
    'https://randomuser.me/api/portraits/women/28.jpg',
    'https://linkedin.com/in/oliviabrown',
    'DesignLoop Agency',
    'America/Los_Angeles',
    'We are working on a project that could use your expertise.',
    false
FROM accounts
WHERE email = 'vishnu.shon@example.com';

INSERT INTO messages (
    account_id,
    name,
    picture,
    profile,
    company,
    timezone,
    message,
    is_sent
)
SELECT
    id,
    'Benjamin Lee',
    'https://randomuser.me/api/portraits/men/37.jpg',
    'https://linkedin.com/in/benjaminlee',
    'TechInnovate',
    'Asia/Singapore',
    NULL,
    false
FROM accounts
WHERE email = 'vishnu.shon@example.com';

SQL