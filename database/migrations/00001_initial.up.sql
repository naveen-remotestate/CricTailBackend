
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
                       user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                       full_name VARCHAR(120) NOT NULL,

                       mobile_number VARCHAR(15) UNIQUE NOT NULL,

                       password_hash TEXT NOT NULL,

                       is_active BOOLEAN DEFAULT TRUE,

                       created_at TIMESTAMP DEFAULT NOW(),

                       updated_at TIMESTAMP DEFAULT NOW()
);

-- USER SESSIONS
CREATE TABLE user_sessions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                               user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

                               archived_at TIMESTAMP,

                               created_at TIMESTAMP DEFAULT NOW()
);



-- TEAMS
CREATE TABLE teams (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                       name VARCHAR(120) UNIQUE NOT NULL,

                       created_by UUID REFERENCES users(user_id),

                       created_at TIMESTAMP DEFAULT NOW(),

                       updated_at TIMESTAMP DEFAULT NOW()
);


-- TEAM PLAYERS
CREATE TABLE team_players (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                              team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,

                              user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

                              batting_position INT,

                              is_captain BOOLEAN DEFAULT FALSE,

                              is_wicket_keeper BOOLEAN DEFAULT FALSE,

                              created_at TIMESTAMP DEFAULT NOW(),

                              UNIQUE(team_id, user_id)
);


-- MATCHES
CREATE TABLE matches (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                         team_a_id UUID NOT NULL REFERENCES teams(id),

                         team_b_id UUID NOT NULL REFERENCES teams(id),

                         toss_winner_team_id UUID REFERENCES teams(id),

                         toss_decision VARCHAR(10)
                             CHECK (toss_decision IN ('BAT', 'BOWL')),

                         batting_first_team_id UUID REFERENCES teams(id),

                         winner_team_id UUID REFERENCES teams(id),

                         overs INT NOT NULL,

                         current_innings_no INT DEFAULT 1,

                         hosted_by UUID REFERENCES users(user_id),

                         scorer_1 UUID REFERENCES users(user_id),

                         scorer_2 UUID REFERENCES users(user_id),

                         stats_processed BOOLEAN DEFAULT FALSE,

                         start_time TIMESTAMP,

                         end_time TIMESTAMP,

                         created_at TIMESTAMP DEFAULT NOW(),

                         updated_at TIMESTAMP DEFAULT NOW()
);


-- INNINGS
CREATE TABLE innings (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                         match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,

                         innings_no INT NOT NULL,

                         batting_team_id UUID NOT NULL REFERENCES teams(id),

                         bowling_team_id UUID NOT NULL REFERENCES teams(id),

                         total_runs INT DEFAULT 0,

                         total_wickets INT DEFAULT 0,

                         legal_balls INT DEFAULT 0,

                         extras INT DEFAULT 0,

                         wides INT DEFAULT 0,

                         no_balls INT DEFAULT 0,

                         byes INT DEFAULT 0,

                         leg_byes INT DEFAULT 0,

                         is_completed BOOLEAN DEFAULT FALSE,

                         start_time TIMESTAMP,

                         end_time TIMESTAMP,

                         created_at TIMESTAMP DEFAULT NOW(),

                         updated_at TIMESTAMP DEFAULT NOW(),

                         UNIQUE(match_id, innings_no)
);


-- LIVE MATCH STATE
CREATE TABLE live_match (
                                  match_id UUID PRIMARY KEY REFERENCES matches(id) ON DELETE CASCADE,

                                  innings_id UUID REFERENCES innings(id),

                                  striker_id UUID REFERENCES users(user_id),

                                  non_striker_id UUID REFERENCES users(user_id),

                                  current_bowler_id UUID REFERENCES users(user_id),

                                  total_runs INT DEFAULT 0,

                                  total_wickets INT DEFAULT 0,

                                  legal_balls INT DEFAULT 0,

                                  updated_at TIMESTAMP DEFAULT NOW()
);


-- BALL EVENTS
CREATE TABLE ball_events (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                             innings_id UUID NOT NULL REFERENCES innings(id) ON DELETE CASCADE,

                             ball_sequence INT NOT NULL,

                             over_no INT NOT NULL,

                             ball_in_over INT NOT NULL,

                             striker_id UUID NOT NULL REFERENCES users(user_id),

                             non_striker_id UUID NOT NULL REFERENCES users(user_id),

                             bowler_id UUID NOT NULL REFERENCES users(user_id),

                             runs_off_bat INT DEFAULT 0,

                             extra_runs INT DEFAULT 0,

                             total_runs INT DEFAULT 0,

                             extra_type VARCHAR(20)
                                 CHECK (
                                     extra_type IN (
                                                    'WIDE',
                                                    'NO_BALL',
                                                    'BYE',
                                                    'LEG_BYE'
                                         )
                                     ),

                             is_legal_delivery BOOLEAN DEFAULT TRUE,

                             is_boundary_four BOOLEAN DEFAULT FALSE,

                             is_boundary_six BOOLEAN DEFAULT FALSE,

                             is_dot_ball BOOLEAN DEFAULT FALSE,

                             is_wicket BOOLEAN DEFAULT FALSE,

                             wicket_type VARCHAR(30)
                                 CHECK (
                                     wicket_type IN (
                                                     'BOWLED',
                                                     'CAUGHT',
                                                     'LBW',
                                                     'RUN_OUT',
                                                     'STUMPED',
                                                     'HIT_WICKET',
                                                     'RETIRED_HURT'
                                         )
                                     ),

                             dismissed_player_id UUID REFERENCES users(user_id),

                             dismissed_by_fielder_id UUID REFERENCES users(user_id),

                             bowled_at TIMESTAMP DEFAULT NOW()
);


-- BATTING SCORECARDS
CREATE TABLE batting_scorecards (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                    innings_id UUID NOT NULL REFERENCES innings(id) ON DELETE CASCADE,

                                    user_id UUID NOT NULL REFERENCES users(user_id),

                                    batting_position INT,

                                    runs INT DEFAULT 0,

                                    balls_faced INT DEFAULT 0,

                                    fours INT DEFAULT 0,

                                    sixes INT DEFAULT 0,

                                    dismissal_type VARCHAR(30)
                                        CHECK (
                                            dismissal_type IN (
                                                               'BOWLED',
                                                               'CAUGHT',
                                                               'LBW',
                                                               'RUN_OUT',
                                                               'STUMPED',
                                                               'HIT_WICKET',
                                                               'RETIRED_HURT'
                                                )
                                            ),

                                    dismissed_by_bowler_id UUID REFERENCES users(user_id),

                                    fielder_id UUID REFERENCES users(user_id),

                                    is_out BOOLEAN DEFAULT FALSE,

                                    created_at TIMESTAMP DEFAULT NOW(),

                                    updated_at TIMESTAMP DEFAULT NOW(),

                                    UNIQUE(innings_id, user_id)
);


-- BOWLING SCORECARDS
CREATE TABLE bowling_scorecards (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                    innings_id UUID NOT NULL REFERENCES innings(id) ON DELETE CASCADE,

                                    user_id UUID NOT NULL REFERENCES users(user_id),

                                    legal_balls INT DEFAULT 0,

                                    maidens INT DEFAULT 0,

                                    runs_conceded INT DEFAULT 0,

                                    wickets INT DEFAULT 0,

                                    wides INT DEFAULT 0,

                                    no_balls INT DEFAULT 0,

                                    created_at TIMESTAMP DEFAULT NOW(),

                                    updated_at TIMESTAMP DEFAULT NOW(),

                                    UNIQUE(innings_id, user_id)
);


-- PLAYER CAREER STATS
CREATE TABLE player_career_stats (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                     user_id UUID UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

                                     batting_style VARCHAR(10)
                                         CHECK (batting_style IN ('LEFT', 'RIGHT')),
                                     bowling_style VARCHAR(10)
                                         CHECK (bowling_style IN ('FAST', 'MEDIUM','SPIN','OFF SPIN','LEG SPIN')),

                                     matches_played INT DEFAULT 0,

                                     innings_batted INT DEFAULT 0,

                                     innings_bowled INT DEFAULT 0,

                                     total_runs INT DEFAULT 0,

                                     highest_score INT DEFAULT 0,

                                     total_balls_faced INT DEFAULT 0,

                                     total_fours INT DEFAULT 0,

                                     total_sixes INT DEFAULT 0,

                                     total_wickets INT DEFAULT 0,

                                     total_balls_bowled INT DEFAULT 0,

                                     total_runs_conceded INT DEFAULT 0,

                                     total_maidens INT DEFAULT 0,

                                     catches INT DEFAULT 0,

                                     run_outs INT DEFAULT 0,

                                     updated_at TIMESTAMP DEFAULT NOW()
);

