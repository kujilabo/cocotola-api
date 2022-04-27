create table `english_phrase_problem` (
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
,`text` varchar(100) character set ascii not null
,`lang2` varchar(2) character set ascii
,`translated` varchar(100)
,primary key(`id`)
,unique(`organization_id`, `workbook_id`, `text`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`updated_by`) references `app_user`(`id`) on delete cascade
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`workbook_id`) references `workbook`(`id`) on delete cascade
);
