package models

import (
	"time"

	"gorm.io/datatypes"
)

type Agendamento struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DataInicio       time.Time  `gorm:"not null" json:"data_inicio"`                // timestamp with time zone
	Numero           *string    `json:"numero"`                                     // text null
	DataFim          *time.Time `json:"data_fim"`                                   // timestamp with time zone null
	IDUser           *uint      `json:"id_user"`                                    // integer null (FK para users)
	NomeCliente      *string    `json:"nome_cliente"`                               // text null
	Indisponivel     bool       `gorm:"not null;default:false" json:"indisponivel"` // boolean not null default false
	IDFuncionario    *uint      `json:"id_funcionario"`                             // integer null (FK para funcionarios)
	EnderecoCliente  *string    `json:"endereco_cliente"`                           // text null
	NotificacaoVista *bool      `gorm:"default:false" json:"notificacao_vista"`     // boolean null default false
	IDPagamento      *uint      `json:"id_pagamento"`                               // integer null (FK para pagamentos)
	DataAniversario  *time.Time `json:"data_aniversario"`                           // date null

	// Foreign keys
	ServicosAgendado []AgendamentoServico `gorm:"foreignKey:IdAgendamento" json:"agendamentos_servicos"`
	Funcionario      *Funcionario         `gorm:"foreignKey:IDFuncionario;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	User             *User                `gorm:"foreignKey:IDUser" json:"-"`
	Pagamento        *Pagamento           `gorm:"foreignKey:IDPagamento;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Funcionario struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Nome   string  `gorm:"not null"` // text not null
	Image  *string // text null
	CPF    *string // text null
	IDUser *uint   `json:"-"` // integer null (FK para users)
	Color  *string // text null

	// Foreign key
	User *User `gorm:"foreignKey:IDUser" json:"-"`
}
type User struct {
	ID               uint           `gorm:"primaryKey;autoIncrement"`
	Email            string         `gorm:"not null"`               // text not null
	Active           bool           `gorm:"not null;default:false"` // boolean not null default false
	Name             *string        // text null
	Senha            *string        // text null
	LogoURL          *string        `gorm:"column:logoUrl"` // text null, com nome especial
	DiasFunciona     datatypes.JSON `gorm:"type:json"`      // integer[] → armazenado como JSON
	Link             *string        // text null
	HasFuncionarios  *bool          `gorm:"default:false"` // boolean null default false
	Contato          *string        // text null
	AtendimentoLocal *bool          `gorm:"column:atendimentoLocal;default:false"` // boolean null default false
	Config           datatypes.JSON `gorm:"type:json"`                             // json null
	Endereco         datatypes.JSON `gorm:"type:json"`                             // json null
	Insta            *string        // text null
	Images           datatypes.JSON `gorm:"type:json"` // text[] → armazenado como JSON
}

type Pagamento struct {
	ID                uint       `gorm:"primaryKey;autoIncrement"`
	PaymentID         *int64     // bigint null
	CollectorID       *int64     `json:"-"` // bigint null
	DateCreated       *time.Time // timestamp without time zone null
	ExternalReference *string    // text null
	PaymentMethodID   *string    // text null
	Status            *string    // text null
	StatusDetail      *string    // text null
	UserID            *uint      // integer null (FK)
	TransactionAmount *float64   // numeric null
	Externo           *bool      // boolean null

	// Foreign key
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

type Servico struct {
	ID          int      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome        string   `gorm:"not null" json:"nome"`
	IdUser      *int     `gorm:"column:id_user" json:"-"`
	Tempo       *string  `json:"tempo"`
	Preco       *float64 `gorm:"type:double precision" json:"preco"`
	Image       *string  `json:"image"`
	Categoria   *string  `json:"categoria"`
	Featured    *bool    `json:"featured"`
	Description *string  `json:"description"`
	Aviso       *string  `json:"aviso"`

	User User `gorm:"foreignKey:IdUser;constraint:OnDelete:CASCADE" json:"-"`
}
type AgendamentoServico struct {
	ID            int `gorm:"primaryKey;autoIncrement" json:"-"`
	IdAgendamento int `gorm:"column:id_agendamento;not null" json:"-"`
	IdServico     int `gorm:"column:id_servicos" json:"-"`

	Agendamento Agendamento `gorm:"foreignKey:IdAgendamento;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Servico     Servico     `gorm:"foreignKey:IdServico;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"servico"`
}

func (AgendamentoServico) TableName() string {
	return "agendamentos_servicos"
}

type Saida struct {
	ID          uint       `gorm:"primaryKey; autoIncrement" json:"id"`
	Date        *time.Time `gorm:"type:date" json:"data"`
	Amount      *float64   `json:"amount"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
	Recorrente  *bool      `json:"recorrente"`
	UserID      *uint      `json:"userid"`
	Type        *string    `json:"type"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
