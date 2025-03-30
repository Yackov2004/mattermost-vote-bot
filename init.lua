-- polls
box.schema.space.create('polls', { if_not_exists = true })
box.space.polls:format({
    { name = 'id',       type = 'unsigned' },
    { name = 'question', type = 'string'   },
    { name = 'options',  type = 'array'    },
    { name = 'active',   type = 'boolean'  },
    { name = 'owner_id', type = 'string'   },
})
box.space.polls:create_index('primary', {
    parts = {1, 'unsigned'},
    if_not_exists = true
})

-- poll_votes
box.schema.space.create('poll_votes', { if_not_exists = true })
box.space.poll_votes:format({
    { name = 'poll_id', type = 'unsigned' },
    { name = 'user_id', type = 'string'   },
    { name = 'option',  type = 'string'   },
})

-- Индекс, чтобы запретить повторное голосование одного пользователя за один poll_id, но разрешить голосовать за разные
box.space.poll_votes:create_index('poll_id_user_id_option', {
    parts = {
        {1, 'unsigned'},
        {2, 'string'},
        {3, 'string'},
    },
    unique = true
})

-- Индекс для поиска по poll_id
box.space.poll_votes:create_index('poll_id', {
    parts = {1, 'unsigned'},
    unique = false,
    if_not_exists = true
})