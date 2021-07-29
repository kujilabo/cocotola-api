create table `english_sentence_problem` (
 `id` integer primary key autoincrement
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`workbook_id` int not null
,`number` int not null
,`text` varchar(100) not null
,`lang` varchar(2) not null
,`translated` varchar(100)
,unique(`organization_id`, `workbook_id`, `text`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
,foreign key(`workbook_id`) references `workbook`(`id`)
);
