create table `group_user` (
 `created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`app_user_group_id` int not null
,`app_user_id` int not null
,unique(`app_user_group_id`, `app_user_id`)
);
