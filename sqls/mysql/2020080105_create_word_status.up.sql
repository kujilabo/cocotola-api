create table `word_status` (
 `id` int auto_increment
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`word` varchar(30) not null
,`translation_status` int not null
,`speech_status` int not null
,`phonetic_status` int not null
,`form_status` int not null
,`tatoeba_status` int not null
,`base_word_status` int not null
,primary key(`id`)
,unique(`word`)
);
