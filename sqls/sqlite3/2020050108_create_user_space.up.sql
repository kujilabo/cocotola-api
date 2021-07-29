create table `user_space` (
 `created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`app_user_id` int not null
,`space_id` int not null
,unique(`app_user_id`, `space_id`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
,foreign key(`app_user_id`) references `app_user`(`id`)
,foreign key(`space_id`) references `space`(`id`)
);
