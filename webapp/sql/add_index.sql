ALTER TABLE livestream_tags ADD INDEX (tag_id);
ALTER TABLE livestream_tags ADD INDEX (livestream_id);
-- ALTER TABLE isudns.records ADD INDEX (domain_id,name);
ALTER TABLE icons ADD INDEX (user_id);
ALTER TABLE isudns.records ADD INDEX (name);
ALTER TABLE themes ADD INDEX (user_id);
