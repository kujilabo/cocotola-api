create table `user_quota` (
 `organization_id` int not null
,`app_user_id` int not null
,`date` datetime not null
,`name` varchar(32) not null
,`unit` varchar(16) not null
,`count` int not null
,unique(`organization_id`, `app_user_id`, `date`, `name`)
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
);
