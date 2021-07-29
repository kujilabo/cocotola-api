create table `user_role` (
 `created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`app_user_id` int not null
,`role` varchar(20) not null
,unique(`app_user_id`, `role`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`app_user_id`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
);
