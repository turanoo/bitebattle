-- Remove users table and all foreign key constraints referencing users(id)
ALTER TABLE polls DROP CONSTRAINT IF EXISTS polls_created_by_fkey;
ALTER TABLE polls_members DROP CONSTRAINT IF EXISTS polls_members_user_id_fkey;
ALTER TABLE head2head_matches DROP CONSTRAINT IF EXISTS head2head_matches_inviter_id_fkey;
ALTER TABLE head2head_matches DROP CONSTRAINT IF EXISTS head2head_matches_invitee_id_fkey;
ALTER TABLE head2head_swipes DROP CONSTRAINT IF EXISTS head2head_swipes_user_id_fkey;
ALTER TABLE votes DROP CONSTRAINT IF EXISTS votes_user_id_fkey;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_user_id_fkey;
DROP TABLE IF EXISTS users;
