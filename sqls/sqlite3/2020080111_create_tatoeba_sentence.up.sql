create table `tatoeba_sentence` (
 `sentence_number` int not null
,`lang` varchar(3) not null
,`text` varchar(500) not null
,`author` varchar(20) not null
,`updated_at` datetime not null
,primary key(`sentence_number`)
);
