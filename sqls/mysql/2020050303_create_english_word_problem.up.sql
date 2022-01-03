create table `english_word_problem` (
 `id` int auto_increment
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`workbook_id` int not null
,`audio_id` int not null
,`number` int not null
,`text` varchar(30) character set ascii not null
,`pos` int not null
,`phonetic` varchar(50)
,`present_third` varchar(30) character set ascii
,`present_participle` varchar(30) character set ascii
,`past_tense` varchar(30) character set ascii
,`past_participle` varchar(30) character set ascii
,`lang` varchar(2) character set ascii
,`translated` varchar(100)
,`phrase_id1` int
,`phrase_id2` int
,`sentence_id1` int
,`sentence_id2` int
,primary key(`id`)
,unique(`organization_id`, `workbook_id`, `text`, `pos`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`updated_by`) references `app_user`(`id`) on delete cascade
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`workbook_id`) references `workbook`(`id`) on delete cascade
);
