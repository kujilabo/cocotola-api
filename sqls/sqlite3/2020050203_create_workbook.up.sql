create table `workbook` (
 `id` integer primary key autoincrement
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`owner_id` int not null
,`space_id` int not null
,`problem_type_id` int not null
,`name` varchar(40) not null
,`question_text` varchar(100)
,unique(`organization_id`, `owner_id`, `name`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
,foreign key(`owner_id`) references `app_user`(`id`)
,foreign key(`space_id`) references `space`(`id`)
,foreign key(`problem_type_id`) references `problem_type`(`id`)
);
