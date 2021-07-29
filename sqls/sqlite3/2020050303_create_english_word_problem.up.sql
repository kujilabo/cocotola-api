create table `english_word_problem` (
 `id` integer primary key autoincrement
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`workbook_id` int not null
,`audio_id` int not null
,`number` int not null
,`text` varchar(30) not null
,`pos` int not null
,`phonetic` varchar(50)
,`present_third` varchar(30)
,`present_participle` varchar(30)
,`past_tense` varchar(30)
,`past_participle` varchar(30)
,`lang` varchar(2)
,`translated` varchar(100)
,`phrase_id1` int
,`phrase_id2` int
,`sentence_id1` int
,`sentence_id2` int
,unique(`organization_id`, `workbook_id`, `text`, `pos`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
,foreign key(`workbook_id`) references `workbook`(`id`)
,foreign key(`audio_id`) references `audio`(`id`)
);
